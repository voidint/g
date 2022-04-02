package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/voidint/g/errs"
)

// FindVersion 返回指定名称的版本
func FindVersion(all []*Version, name string) (*Version, error) {
	for i := range all {
		if all[i].Name == name {
			return all[i], nil
		}
	}
	return nil, errs.ErrVersionNotFound
}

// Version go版本
type Version struct {
	Name     string // 版本名，如'1.12.4'
	Packages []*Package
}

// FindPackage 返回指定操作系统和硬件架构的版本包
func (v *Version) FindPackage(kind, goos, goarch string) (*Package, error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		return v.Packages[i], nil
	}

	return nil, errs.ErrPackageNotFound
}

// FindPackages 返回指定操作系统和硬件架构的版本包
func (v *Version) FindPackages(kind, goos, goarch string) (pkgs []*Package, err error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		pkgs = append(pkgs, v.Packages[i])
	}
	if len(pkgs) == 0 {
		return nil, errs.ErrPackageNotFound
	}
	return pkgs, nil
}

// Package go版本安装包
type Package struct {
	FileName  string
	URL       string
	Kind      string
	OS        string
	Arch      string
	Size      string
	Checksum  string
	Algorithm string // checksum algorithm
}

const (
	// SourceKind go安装包种类-源码
	SourceKind = "Source"
	// ArchiveKind go安装包种类-压缩文件
	ArchiveKind = "Archive"
	// InstallerKind go安装包种类-可安装程序
	InstallerKind = "Installer"
)

// Download 下载版本另存为指定文件
func (pkg *Package) Download(dst string) (size int64, err error) {
	resp, err := http.Get(pkg.URL)
	if err != nil {
		return 0, errs.NewDownloadError(pkg.URL, err)
	}
	defer resp.Body.Close()
	f, err := os.Create(dst)
	if err != nil {
		return 0, errs.NewDownloadError(pkg.URL, err)
	}
	defer f.Close()
	size, err = io.Copy(f, resp.Body)
	if err != nil {
		return 0, errs.NewDownloadError(pkg.URL, err)
	}
	return size, nil
}

// DownloadWithProgress 下载版本另存为指定文件且显示下载进度
func (pkg *Package) DownloadWithProgress(dst string) (size int64, err error) {
	resp, err := http.Get(pkg.URL)
	if err != nil {
		return 0, errs.NewDownloadError(pkg.URL, err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, errs.NewDownloadError(pkg.URL, err)
	}
	defer f.Close()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(ansi.NewAnsiStdout(), "\n")
		}),
		// progressbar.OptionSpinnerType(35),
		// progressbar.OptionFullWidth(),
	)
	_ = bar.RenderBlank()

	size, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return size, errs.NewDownloadError(pkg.URL, err)
	}
	return size, nil
}

const (
	// SHA256 校验和算法-sha256
	SHA256 = "SHA256"
	// SHA1 校验和算法-sha1
	SHA1 = "SHA1"
)

// VerifyChecksum 验证目标文件的校验和与当前安装包的校验和是否一致
func (pkg *Package) VerifyChecksum(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	var h hash.Hash
	switch pkg.Algorithm {
	case SHA256:
		h = sha256.New()
	case SHA1:
		h = sha1.New()
	default:
		return errs.ErrUnsupportedChecksumAlgorithm
	}

	if _, err = io.Copy(h, f); err != nil {
		return err
	}

	if pkg.Checksum != hex.EncodeToString(h.Sum(nil)) {
		return errs.ErrChecksumNotMatched
	}
	return nil
}
