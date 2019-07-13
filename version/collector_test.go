package version

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/PuerkitoBio/goquery"

	. "github.com/smartystreets/goconvey/convey"
)

func getCollector() (*Collector, error) {
	b, err := ioutil.ReadFile("./testdata/golang_dl.html")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &Collector{
		url: DefaultURL,
		doc: doc,
	}, nil
}

func Test_findPackages(t *testing.T) {
	Convey("查找目标go版本下的安装包列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		pkgs := c.findPackages(c.doc.Find("#unstable").Next().Find("table").First())
		So(len(pkgs), ShouldEqual, 15)

		for i, expected := range []*Package{
			{
				FileName:  "go1.13beta1.src.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.src.tar.gz",
				Kind:      SourceKind,
				OS:        "",
				Arch:      "",
				Size:      "21MB",
				Checksum:  "e8a7c504cd6775b8a6af101158b8871455918c9a61162f0180f7a9f118dc4102",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.darwin-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.darwin-amd64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "macOS",
				Arch:      "x86-64",
				Size:      "117MB",
				Checksum:  "7af1aead60905c14085300b38a39b8ea2da5d6bf55084caa759a8bdf41ae0c32",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.darwin-amd64.pkg",
				URL:       "https://dl.google.com/go/go1.13beta1.darwin-amd64.pkg",
				Kind:      InstallerKind,
				OS:        "macOS",
				Arch:      "x86-64",
				Size:      "116MB",
				Checksum:  "f7f0a0dd1fb18337e182fc0d93ecc71622b36fb3dfa2644a4f8bc0f67aa5f84d",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.freebsd-386.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.freebsd-386.tar.gz",
				Kind:      ArchiveKind,
				OS:        "FreeBSD",
				Arch:      "x86",
				Size:      "96MB",
				Checksum:  "b9505fa721ab1e8c972172374fa2db52e67955798c5c8574620f74bd7900a808",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.freebsd-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.freebsd-amd64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "FreeBSD",
				Arch:      "x86-64",
				Size:      "114MB",
				Checksum:  "9c1fb2edaf403bba04d49f2f7da4d09b14c63bbe6143f1ff1e8ba56b4e17d013",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-386.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-386.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "x86",
				Size:      "97MB",
				Checksum:  "38039e4f7b6eea8f55e91d90607150d5d397f9063c06445c45009dd1e6dba8cc",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-amd64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-amd64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "x86-64",
				Size:      "114MB",
				Checksum:  "dbd131c92f381a5bc5ca1f0cfd942cb8be7d537007b6f412b5be41ff38a7d0d9",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-arm64.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-arm64.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "ARMv8",
				Size:      "103MB",
				Checksum:  "298a325d8eeba561a26312a9cdc821a96873c10fca7f48a7f98bbd8848bd8bd4",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-armv6l.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-armv6l.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "ARMv6",
				Size:      "94MB",
				Checksum:  "77993f1dce5b4d080cbd06a4553e5e1c6caa7ad6817ea3c62254b89d6f079504",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-ppc64le.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-ppc64le.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "ppc64le",
				Size:      "92MB",
				Checksum:  "0f3c5c7b7956911ed8d1fc4e9dbeb2584d0be695c5c15b528422e3bb2d5989f0",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.linux-s390x.tar.gz",
				URL:       "https://dl.google.com/go/go1.13beta1.linux-s390x.tar.gz",
				Kind:      ArchiveKind,
				OS:        "Linux",
				Arch:      "s390x",
				Size:      "97MB",
				Checksum:  "877065ac7d1729e5de1bbfe1e712788bf9dee5613a5502cf0ba76e65c2521b26",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.windows-386.zip",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-386.zip",
				Kind:      ArchiveKind,
				OS:        "Windows",
				Arch:      "x86",
				Size:      "111MB",
				Checksum:  "f0908f1703c642950442317f7581c8254842f00298e4e0f511d1513c87e3c64d",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.windows-386.msi",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-386.msi",
				Kind:      InstallerKind,
				OS:        "Windows",
				Arch:      "x86",
				Size:      "96MB",
				Checksum:  "6189e5d13ef054117fc45fe028a4b3c6b22fc8301a422e6fb13f332a864a8da9",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.windows-amd64.zip",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-amd64.zip",
				Kind:      ArchiveKind,
				OS:        "Windows",
				Arch:      "x86-64",
				Size:      "129MB",
				Checksum:  "08098b4b0e1a105971d2fced2842e806f8ffa08973ae8781fd22dd90f76404fb",
				Algorithm: SHA256,
			},
			{
				FileName:  "go1.13beta1.windows-amd64.msi",
				URL:       "https://dl.google.com/go/go1.13beta1.windows-amd64.msi",
				Kind:      InstallerKind,
				OS:        "Windows",
				Arch:      "x86-64",
				Size:      "112MB",
				Checksum:  "989098d4f3535ebd0126f381eb9ff097373c060ad8fce902730696866e84f297",
				Algorithm: SHA256,
			},
		} {
			So(pkgs[i].Algorithm, ShouldEqual, expected.Algorithm)
			So(pkgs[i].FileName, ShouldEqual, expected.FileName)
			So(pkgs[i].Kind, ShouldEqual, expected.Kind)
			So(pkgs[i].OS, ShouldEqual, expected.OS)
			So(pkgs[i].Arch, ShouldEqual, expected.Arch)
			So(pkgs[i].Size, ShouldEqual, expected.Size)
			So(pkgs[i].Checksum, ShouldEqual, expected.Checksum)
		}
	})
}

func TestStableVersions(t *testing.T) {
	Convey("查询stable状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.StableVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 2)

		for i, expected := range []struct {
			Name          string
			PackageLength int
		}{
			{Name: "1.12.7", PackageLength: 15},
			{Name: "1.11.12", PackageLength: 15},
		} {
			So(items[i].Name, ShouldEqual, expected.Name)
			So(len(items[i].Packages), ShouldEqual, expected.PackageLength)
		}
	})
}

func TestUnstableVersions(t *testing.T) {
	Convey("查询unstable状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.UnstableVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 1)
		So(items[0].Name, ShouldEqual, "1.13beta1")
		So(len(items[0].Packages), ShouldEqual, 15)
	})
}

func TestArchivedVersions(t *testing.T) {
	Convey("查询archived状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.ArchivedVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 70)

		So(items[0].Name, ShouldEqual, "1.12.6")
		So(len(items[0].Packages), ShouldEqual, 15)
	})
}

func TestAllVersions(t *testing.T) {
	Convey("查询所有go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.AllVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 73)

		So(items[len(items)-1].Name, ShouldEqual, "1.13beta1")
		So(len(items[len(items)-1].Packages), ShouldEqual, 15)
	})
}

func TestURLUnreachableError(t *testing.T) {
	Convey("URL不可达错误", t, func() {
		url := "https://github.com/voidint"
		core := errors.New("hello error")

		err := NewURLUnreachableError(url, core)
		So(err, ShouldNotBeNil)
		e, ok := err.(*URLUnreachableError)
		So(ok, ShouldBeTrue)
		So(e, ShouldNotBeNil)
		So(e.url, ShouldEqual, url)
		So(e.err, ShouldEqual, core)
		So(e.Error(), ShouldEqual, fmt.Sprintf("URL %q is unreachable ==> %s", url, core.Error()))
	})
}
