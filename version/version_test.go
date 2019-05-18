package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFindVersion(t *testing.T) {
	Convey("查找指定名称的版本", t, func() {
		v0 := &Version{
			Name: "1.12.5",
		}
		v1 := &Version{
			Name: "1.11.10",
		}
		v2 := &Version{
			Name: "1.9.7",
		}

		items := []*Version{v0, v1, v2}

		v, err := FindVersion(items, "1.11.10")
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.Name, ShouldEqual, "1.11.10")

		v, err = FindVersion(items, "1.11.11")
		So(err, ShouldEqual, ErrVersionNotFound)
		So(v, ShouldBeNil)
	})
}

func TestFindPackage(t *testing.T) {
	Convey("查询版本下的安装包", t, func() {
		v := &Version{
			Name: "1.12.4",
			Packages: []*Package{
				{
					FileName: "go1.12.4.src.tar.gz",
					Kind:     SourceKind,
					Size:     "21MB",
				},
				{
					FileName: "go1.12.4.darwin-amd64.tar.gz",
					Kind:     ArchiveKind,
					OS:       "macOS",
					Arch:     "x86-64",
					Size:     "122MB",
				},
				{
					FileName: "go1.12.4.windows-386.msi",
					Kind:     InstallerKind,
					OS:       "Windows",
					Arch:     "x86",
					Size:     "102MB",
				},
			},
		}

		pkg, err := v.FindPackage(ArchiveKind, "darwin", "amd64")
		So(err, ShouldBeNil)
		So(pkg, ShouldNotBeNil)
		So(pkg.FileName, ShouldEqual, "go1.12.4.darwin-amd64.tar.gz")
		So(pkg.Kind, ShouldEqual, ArchiveKind)
		So(pkg.OS, ShouldEqual, "macOS")
		So(pkg.Arch, ShouldEqual, "x86-64")

		pkg, err = v.FindPackage(ArchiveKind, "darwin", "386")
		So(err, ShouldEqual, ErrPackageNotFound)
		So(pkg, ShouldBeNil)
	})
}

func TestDownloadError(t *testing.T) {
	Convey("安装包下载错误", t, func() {
		url := "https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz"
		core := errors.New("hello error")

		err := NewDownloadError(url, core)
		So(err, ShouldNotBeNil)
		e, ok := err.(*DownloadError)
		So(ok, ShouldBeTrue)
		So(e, ShouldNotBeNil)
		So(e.url, ShouldEqual, url)
		So(e.err, ShouldEqual, core)
		So(e.Error(), ShouldEqual, fmt.Sprintf("Installation package(%s) download failed ==> %s", url, core.Error()))
	})
}

func TestVerifyChecksum(t *testing.T) {
	Convey("检查安装包校验和", t, func() {
		filename := fmt.Sprintf("%d.txt", time.Now().Unix())
		f, err := os.Create(filename)
		So(err, ShouldBeNil)
		defer os.Remove(filename)
		defer f.Close()
		_, err = f.WriteString("hello 世界！")
		So(err, ShouldBeNil)

		Convey("SHA256", func() {
			f.Seek(0, 0)
			h := sha256.New()
			_, err = io.Copy(h, f)
			So(err, ShouldBeNil)

			pkg := &Package{
				Algorithm: "SHA256",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			So(pkg.VerifyChecksum(filename), ShouldBeNil)
		})

		Convey("SHA1", func() {
			f.Seek(0, 0)
			h := sha1.New()
			_, err = io.Copy(h, f)
			So(err, ShouldBeNil)

			pkg := &Package{
				Algorithm: "SHA1",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			So(pkg.VerifyChecksum(filename), ShouldBeNil)
		})

		Convey("SHA1024", func() {
			pkg := &Package{
				Algorithm: "SHA1024",
			}
			So(pkg.VerifyChecksum(filename), ShouldEqual, ErrUnsupportedChecksumAlgorithm)
		})
	})
}
