package fancyindex

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/collector/internal"
	"github.com/voidint/g/pkg/errs"
	"github.com/voidint/g/version"
)

const AliYunDownloadPageURL = "https://mirrors.aliyun.com/golang/"

func getCollector() (*Collector, error) {
	b, err := os.ReadFile("./testdata/aliyun.html")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &Collector{
		url: AliYunDownloadPageURL,
		doc: doc,
	}, nil
}

func Test_findGoFileItems(t *testing.T) {
	c, err := getCollector()
	assert.Nil(t, err)
	assert.NotNil(t, c)

	t.Run("", func(t *testing.T) {
		items := c.findGoFileItems(c.doc.Find(".table"))
		assert.True(t, len(items) >= 11)

		for i, gfi := range []*internal.GoFileItem{
			{FileName: "go1.10.1.darwin-amd64.pkg", URL: AliYunDownloadPageURL + "go1.10.1.darwin-amd64.pkg", Size: "111.5 MB"},
			{FileName: "go1.10.1.darwin-amd64.pkg.sha256", URL: AliYunDownloadPageURL + "go1.10.1.darwin-amd64.pkg.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz", URL: AliYunDownloadPageURL + "go1.10.1.darwin-amd64.tar.gz", Size: "112.4 MB"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz.asc", URL: AliYunDownloadPageURL + "go1.10.1.darwin-amd64.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz.sha256", URL: AliYunDownloadPageURL + "go1.10.1.darwin-amd64.tar.gz.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.freebsd-386.tar.gz", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-386.tar.gz", Size: "99.0 MB"},
			{FileName: "go1.10.1.freebsd-386.tar.gz.asc", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-386.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.freebsd-386.tar.gz.sha256", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-386.tar.gz.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz", Size: "110.3 MB"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz.asc", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz.sha256", URL: AliYunDownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz.sha256", Size: "64.0 B"},
		} {
			assert.Equal(t, gfi.FileName, items[i].FileName)
			assert.Equal(t, gfi.URL, items[i].URL)
			assert.Equal(t, gfi.Size, items[i].Size)
		}
	})
}

func TestCollector_StableVersions(t *testing.T) {
	t.Run("稳定版本列表", func(t *testing.T) {
		c := &Collector{}
		vs, err := c.StableVersions()
		assert.Nil(t, err)
		assert.Equal(t, []*version.Version{}, vs)
	})
}

func TestCollector_UnstableVersions(t *testing.T) {
	t.Run("非稳定版本列表", func(t *testing.T) {
		c := &Collector{}
		vs, err := c.UnstableVersions()
		assert.Nil(t, err)
		assert.Equal(t, []*version.Version{}, vs)
	})
}

func TestCollector_ArchivedVersions(t *testing.T) {
	t.Run("已归档版本列表", func(t *testing.T) {
		c := &Collector{}
		vs, err := c.ArchivedVersions()
		assert.Nil(t, err)
		assert.Equal(t, []*version.Version{}, vs)
	})
}

func TestNewCollector(t *testing.T) {
	t.Run("空URL", func(t *testing.T) {
		c, err := NewCollector("")
		assert.Equal(t, errs.ErrEmptyURL, err)
		assert.Nil(t, c)
	})

	rr1 := httptest.NewRecorder()
	rr1.WriteHeader(http.StatusNotFound)

	rr2 := httptest.NewRecorder()
	rr2.WriteHeader(http.StatusOK)
	htmlData, err := os.ReadFile("./testdata/aliyun.html")
	assert.Nil(t, err)
	_, _ = rr2.Write(htmlData)

	patches := gomonkey.ApplyMethodSeq(&http.Client{}, "Get", []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("unknown error")}},
		{Values: gomonkey.Params{rr1.Result(), nil}},
		{Values: gomonkey.Params{rr2.Result(), nil}},
	})
	defer patches.Reset()

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "站点URL访问异常",
			wantErr: errs.NewURLUnreachableError(AliYunDownloadPageURL, errors.New("unknown error")),
		},
		{
			name:    "站点URL资源不存在",
			wantErr: errs.NewURLUnreachableError(AliYunDownloadPageURL, fmt.Errorf("%d", http.StatusNotFound)),
		},
		{
			name:    "站点URL访问采集正常",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollector(AliYunDownloadPageURL)
			assert.Equal(t, tt.wantErr, err)
			if err == nil {
				assert.NotNil(t, got.pURL)
				assert.NotNil(t, got.doc)
				// assert.NotNil(t, got.(*Collector).pURL)
				// assert.NotNil(t, got.(*Collector).doc)
			}
		})
	}
}

func TestCollector_Name(t *testing.T) {
	t.Run("Collector name", func(t *testing.T) {
		c := &Collector{}
		assert.Equal(t, Name, c.Name())
	})
}
