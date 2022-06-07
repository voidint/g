package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/pkg/errs"
)

func TestFindVersion(t *testing.T) {
	t.Run("查找指定名称的版本", func(t *testing.T) {
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
		assert.Nil(t, err)
		assert.NotNil(t, v)
		assert.Equal(t, "1.11.10", v.Name)

		v, err = FindVersion(items, "1.11.11")
		assert.Equal(t, errs.ErrVersionNotFound, err)
		assert.Nil(t, v)
	})
}

func TestFindPackage(t *testing.T) {
	t.Run("查询版本下的安装包", func(t *testing.T) {
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
		assert.Nil(t, err)
		assert.NotNil(t, pkg)
		assert.Equal(t, "go1.12.4.darwin-amd64.tar.gz", pkg.FileName)
		assert.Equal(t, ArchiveKind, pkg.Kind)
		assert.Equal(t, "macOS", pkg.OS)
		assert.Equal(t, "x86-64", pkg.Arch)

		pkg, err = v.FindPackage(ArchiveKind, "darwin", "386")
		assert.Equal(t, errs.ErrPackageNotFound, err)
		assert.Nil(t, pkg)
	})
}

func TestVerifyChecksum(t *testing.T) {
	t.Run("检查安装包校验和", func(t *testing.T) {
		filename := fmt.Sprintf("%d.txt", time.Now().Unix())
		f, err := os.Create(filename)
		assert.Nil(t, err)
		defer os.Remove(filename)
		defer f.Close()
		_, err = f.WriteString("hello 世界！")
		assert.Nil(t, err)

		t.Run("SHA256", func(t *testing.T) {
			_, _ = f.Seek(0, 0)
			h := sha256.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA256",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			assert.Nil(t, pkg.VerifyChecksum(filename))
		})

		t.Run("校验和不匹配", func(t *testing.T) {
			f.Seek(0, 0)
			h := sha1.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA1",
				Checksum:  fmt.Sprintf("%x", h.Sum(nil)),
			}
			assert.Nil(t, pkg.VerifyChecksum(filename))
		})

		t.Run("SHA1", func(t *testing.T) {
			f.Seek(0, 0)
			h := sha1.New()
			_, err = io.Copy(h, f)
			assert.Nil(t, err)

			pkg := &Package{
				Algorithm: "SHA1",
				Checksum:  "hello",
			}
			assert.Equal(t, errs.ErrChecksumNotMatched, pkg.VerifyChecksum(filename))
		})

		t.Run("SHA1024", func(t *testing.T) {
			pkg := &Package{
				Algorithm: "SHA1024",
			}
			assert.Equal(t, errs.ErrUnsupportedChecksumAlgorithm, pkg.VerifyChecksum(filename))
		})
	})
}
