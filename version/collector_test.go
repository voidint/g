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

		pkgs := c.findPackages(c.doc.Find("#stable").Next().Find("table").First())
		So(len(pkgs), ShouldEqual, 15)
		So(pkgs[1].Algorithm, ShouldEqual, "SHA256")
		So(pkgs[1].FileName, ShouldEqual, "go1.12.4.darwin-amd64.tar.gz")
		So(pkgs[1].Kind, ShouldEqual, ArchiveKind)
		So(pkgs[1].OS, ShouldEqual, "macOS")
		So(pkgs[1].Arch, ShouldEqual, "x86-64")
		So(pkgs[1].Size, ShouldEqual, "122MB")
		So(pkgs[1].Checksum, ShouldEqual, "50af1aa6bf783358d68e125c5a72a1ba41fb83cee8f25b58ce59138896730a49")
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
		So(items[0].Name, ShouldEqual, "1.12.4")
		So(len(items[0].Packages), ShouldEqual, 15)
		So(items[1].Name, ShouldEqual, "1.11.9")
		So(len(items[1].Packages), ShouldEqual, 15)
	})
}

func TestArchivedVersions(t *testing.T) {
	Convey("查询archived状态的go版本列表", t, func() {
		c, err := getCollector()
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		items, err := c.ArchivedVersions()
		So(err, ShouldBeNil)
		So(len(items), ShouldEqual, 64)

		So(items[0].Name, ShouldEqual, "1.12.3")
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
		So(len(items), ShouldEqual, 66)
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
