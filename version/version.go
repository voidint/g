package version

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
)

// ErrVersionNotFound 版本不存在
var ErrVersionNotFound = errors.New("version not found")

// FindVersion 返回指定名称的版本
func FindVersion(all []*Version, name string) (*Version, error) {
	for i := range all {
		if all[i].Name == name {
			return all[i], nil
		}
	}
	return nil, ErrVersionNotFound
}

// Version go版本
type Version struct {
	Name     string // 版本名，如'1.12.4'
	Packages []*Package
}

// ErrPackageNotFound 版本包不存在
var ErrPackageNotFound = errors.New("installation package not found")

// FindPackage 返回指定操作系统和硬件架构的版本包
func (v *Version) FindPackage(kind, goos, goarch string) (*Package, error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.Name, goos, goarch)
	for i := range v.Packages {
		if v.Packages[i] == nil || !strings.EqualFold(v.Packages[i].Kind, kind) || !strings.HasPrefix(v.Packages[i].FileName, prefix) {
			continue
		}
		return v.Packages[i], nil
	}
	return nil, ErrPackageNotFound
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

// Download 下载版本另存为指定文件并校验sha256哈希值
func (pkg *Package) Download(dst string) (size int64, err error) {
	resp, err := http.Get(pkg.URL)
	if err != nil {
		return 0, NewDownloadError(pkg.URL, err)
	}
	defer resp.Body.Close()
	f, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	size, err = io.Copy(f, resp.Body)
	if err != nil {
		return 0, NewDownloadError(pkg.URL, err)
	}
	return size, nil
}

// DownloadError 下载失败错误
type DownloadError struct {
	url string
	err error
}

// NewDownloadError 返回下载失败错误实例
func NewDownloadError(url string, err error) error {
	return &DownloadError{
		url: url,
		err: err,
	}
}

func (e *DownloadError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Installation package(%s) download failed", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

var (
	// ErrUnsupportedChecksumAlgorithm 不支持的校验和算法
	ErrUnsupportedChecksumAlgorithm = errors.New("unsupported checksum algorithm")
	// ErrChecksumNotMatched 校验和不匹配
	ErrChecksumNotMatched = errors.New("file checksum does not match the computed checksum")
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
	case "SHA256":
		h = sha256.New()
	case "SHA1":
		h = sha1.New()
	default:
		return ErrUnsupportedChecksumAlgorithm
	}

	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if pkg.Checksum != fmt.Sprintf("%x", h.Sum(nil)) {
		return ErrChecksumNotMatched
	}
	return nil
}
