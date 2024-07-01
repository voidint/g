package official

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/pkg/checksum"
	"github.com/voidint/g/pkg/errs"
	"github.com/voidint/g/version"
)

const OfficialDownloadPageURL = "https://go.dev/dl/"

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
		url: OfficialDownloadPageURL,
		doc: doc,
	}, nil
}

func Test_findPackages(t *testing.T) {
	t.Run("查找目标go版本下的安装包列表", func(t *testing.T) {
		c, err := getCollector()
		assert.Nil(t, err)
		assert.NotNil(t, c)

		pkgs := c.findPackages(c.doc.Find("#unstable").Next().Find("table").First())
		assert.Equal(t, 15, len(pkgs))

		for i, expected := range []*version.Package{
			{
				FileName:  "go1.13beta1.src.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.src.tar.gz",
				Kind:      version.SourceKind,
				OS:        "",
				Arch:      "",
				Size:      "21MB",
				Checksum:  "e8a7c504cd6775b8a6af101158b8871455918c9a61162f0180f7a9f118dc4102",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.darwin-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.darwin-amd64.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "macOS",
				Arch:      "x86-64",
				Size:      "117MB",
				Checksum:  "7af1aead60905c14085300b38a39b8ea2da5d6bf55084caa759a8bdf41ae0c32",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.darwin-amd64.pkg",
				URL:       "https://dl.google.com/go/go1.13beta1.darwin-amd64.pkg",
				Kind:      version.InstallerKind,
				OS:        "macOS",
				Arch:      "x86-64",
				Size:      "116MB",
				Checksum:  "f7f0a0dd1fb18337e182fc0d93ecc71622b36fb3dfa2644a4f8bc0f67aa5f84d",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.freebsd-386.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.freebsd-386.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "FreeBSD",
				Arch:      "x86",
				Size:      "96MB",
				Checksum:  "b9505fa721ab1e8c972172374fa2db52e67955798c5c8574620f74bd7900a808",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.freebsd-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.freebsd-amd64.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "FreeBSD",
				Arch:      "x86-64",
				Size:      "114MB",
				Checksum:  "9c1fb2edaf403bba04d49f2f7da4d09b14c63bbe6143f1ff1e8ba56b4e17d013",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-386.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-386.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "x86",
				Size:      "97MB",
				Checksum:  "38039e4f7b6eea8f55e91d90607150d5d397f9063c06445c45009dd1e6dba8cc",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-amd64.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "x86-64",
				Size:      "114MB",
				Checksum:  "dbd131c92f381a5bc5ca1f0cfd942cb8be7d537007b6f412b5be41ff38a7d0d9",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-arm64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-arm64.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "ARMv8",
				Size:      "103MB",
				Checksum:  "298a325d8eeba561a26312a9cdc821a96873c10fca7f48a7f98bbd8848bd8bd4",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-armv6l.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-armv6l.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "ARMv6",
				Size:      "94MB",
				Checksum:  "77993f1dce5b4d080cbd06a4553e5e1c6caa7ad6817ea3c62254b89d6f079504",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-ppc64le.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-ppc64le.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "ppc64le",
				Size:      "92MB",
				Checksum:  "0f3c5c7b7956911ed8d1fc4e9dbeb2584d0be695c5c15b528422e3bb2d5989f0",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.linux-s390x.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-s390x.tar.gz",
				Kind:      version.ArchiveKind,
				OS:        "Linux",
				Arch:      "s390x",
				Size:      "97MB",
				Checksum:  "877065ac7d1729e5de1bbfe1e712788bf9dee5613a5502cf0ba76e65c2521b26",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.windows-386.zip",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-386.zip",
				Kind:      version.ArchiveKind,
				OS:        "Windows",
				Arch:      "x86",
				Size:      "111MB",
				Checksum:  "f0908f1703c642950442317f7581c8254842f00298e4e0f511d1513c87e3c64d",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.windows-386.msi",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-386.msi",
				Kind:      version.InstallerKind,
				OS:        "Windows",
				Arch:      "x86",
				Size:      "96MB",
				Checksum:  "6189e5d13ef054117fc45fe028a4b3c6b22fc8301a422e6fb13f332a864a8da9",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.windows-amd64.zip",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-amd64.zip",
				Kind:      version.ArchiveKind,
				OS:        "Windows",
				Arch:      "x86-64",
				Size:      "129MB",
				Checksum:  "08098b4b0e1a105971d2fced2842e806f8ffa08973ae8781fd22dd90f76404fb",
				Algorithm: string(checksum.SHA256),
			},
			{
				FileName:  "go1.13beta1.windows-amd64.msi",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-amd64.msi",
				Kind:      version.InstallerKind,
				OS:        "Windows",
				Arch:      "x86-64",
				Size:      "112MB",
				Checksum:  "989098d4f3535ebd0126f381eb9ff097373c060ad8fce902730696866e84f297",
				Algorithm: string(checksum.SHA256),
			},
		} {
			assert.Equal(t, expected.Algorithm, pkgs[i].Algorithm)
			assert.Equal(t, expected.FileName, pkgs[i].FileName)
			assert.Equal(t, expected.Kind, pkgs[i].Kind)
			assert.Equal(t, expected.OS, pkgs[i].OS)
			assert.Equal(t, expected.Arch, pkgs[i].Arch)
			assert.Equal(t, expected.Size, pkgs[i].Size)
			assert.Equal(t, expected.Checksum, pkgs[i].Checksum)
		}
	})
}

func TestUnstableVersions(t *testing.T) {
	t.Run("查询unstable状态的go版本列表", func(t *testing.T) {
		c, err := getCollector()
		assert.Nil(t, err)
		assert.NotNil(t, c)

		items, err := c.UnstableVersions()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(items))
		assert.Equal(t, "1.13beta1", items[0].Name())
		assert.Equal(t, 15, len(items[0].Packages()))
	})
}

func TestArchivedVersions(t *testing.T) {
	t.Run("查询archived状态的go版本列表", func(t *testing.T) {
		c, err := getCollector()
		assert.Nil(t, err)
		assert.NotNil(t, c)

		items, err := c.ArchivedVersions()
		assert.Nil(t, err)
		assert.Equal(t, 70, len(items))
		assert.Equal(t, "1", items[0].Name())
		assert.Equal(t, 2, len(items[0].Packages()))
		assert.Equal(t, "1.12.6", items[len(items)-1].Name())
		assert.Equal(t, 15, len(items[len(items)-1].Packages()))
	})
}

func TestAllVersions(t *testing.T) {
	t.Run("查询所有go版本列表", func(t *testing.T) {
		c, err := getCollector()
		assert.Nil(t, err)
		assert.NotNil(t, c)

		items, err := c.AllVersions()
		assert.Nil(t, err)
		assert.Equal(t, 73, len(items))
		assert.Equal(t, "1.13beta1", items[len(items)-1].Name())
		assert.Equal(t, 15, len(items[len(items)-1].Packages()))
	})
}

func TestNewCollector(t *testing.T) {
	t.Run("空URL", func(t *testing.T) {
		c, err := NewCollector("")
		assert.Equal(t, errs.ErrEmptyURL, err)
		assert.Nil(t, c)
	})

	t.Run("无效URL", func(t *testing.T) {
		var invalidURL strings.Builder
		invalidURL.WriteByte(0x7f)
		invalidURL.WriteString("hello world")

		c, err := NewCollector(invalidURL.String())
		assert.Nil(t, c)
		assert.NotNil(t, err)
		e, ok := err.(*url.Error)
		assert.True(t, ok)
		assert.Equal(t, "parse", e.Op)
		assert.Equal(t, invalidURL.String(), e.URL)
		assert.NotNil(t, e.Err)
	})

	rr1 := httptest.NewRecorder()
	rr1.WriteHeader(http.StatusNotFound)

	rr2 := httptest.NewRecorder()
	rr2.WriteHeader(http.StatusOK)
	htmlData, err := os.ReadFile("./testdata/golang_dl_with_rc.html")
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
		url     string
		wantErr error
	}{
		{
			name:    "默认站点URL访问异常",
			url:     OfficialDownloadPageURL,
			wantErr: errs.NewURLUnreachableError(OfficialDownloadPageURL, errors.New("unknown error")),
		},
		{
			name:    "默认站点URL资源不存在",
			url:     OfficialDownloadPageURL,
			wantErr: errs.NewURLUnreachableError(OfficialDownloadPageURL, fmt.Errorf("%d", http.StatusNotFound)),
		},
		{
			name:    "默认站点URL访问采集正常",
			url:     OfficialDownloadPageURL,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCollector(tt.url)
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
