package aliyun

import (
	"bytes"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/version"
)

func getCollector() (*Collector, error) {
	b, err := os.ReadFile("./testdata/golang_dl.html")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &Collector{
		url: DownloadPageURL,
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

		for i, gfi := range []*goFileItem{
			{FileName: "go1.10.1.darwin-amd64.pkg", URL: DownloadPageURL + "go1.10.1.darwin-amd64.pkg", Size: "111.5 MB"},
			{FileName: "go1.10.1.darwin-amd64.pkg.sha256", URL: DownloadPageURL + "go1.10.1.darwin-amd64.pkg.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz", URL: DownloadPageURL + "go1.10.1.darwin-amd64.tar.gz", Size: "112.4 MB"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz.asc", URL: DownloadPageURL + "go1.10.1.darwin-amd64.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.darwin-amd64.tar.gz.sha256", URL: DownloadPageURL + "go1.10.1.darwin-amd64.tar.gz.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.freebsd-386.tar.gz", URL: DownloadPageURL + "go1.10.1.freebsd-386.tar.gz", Size: "99.0 MB"},
			{FileName: "go1.10.1.freebsd-386.tar.gz.asc", URL: DownloadPageURL + "go1.10.1.freebsd-386.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.freebsd-386.tar.gz.sha256", URL: DownloadPageURL + "go1.10.1.freebsd-386.tar.gz.sha256", Size: "64.0 B"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz", URL: DownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz", Size: "110.3 MB"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz.asc", URL: DownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz.asc", Size: "819.0 B"},
			{FileName: "go1.10.1.freebsd-amd64.tar.gz.sha256", URL: DownloadPageURL + "go1.10.1.freebsd-amd64.tar.gz.sha256", Size: "64.0 B"},
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
