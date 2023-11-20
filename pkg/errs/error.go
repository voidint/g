package errs

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrUnsupportedChecksumAlgorithm 不支持的校验和算法
	ErrUnsupportedChecksumAlgorithm = errors.New("unsupported checksum algorithm")
	// ErrChecksumNotMatched 校验和不匹配
	ErrChecksumNotMatched = errors.New("file checksum does not match the computed checksum")
	// ErrChecksumFileNotFound 校验和文件不存在
	ErrChecksumFileNotFound = errors.New("checksum file not found")
)

// PackageNotFoundError 软件包不存在错误
type PackageNotFoundError struct {
	kind   string
	goos   string
	goarch string
}

func NewPackageNotFoundError(kind, goos, goarch string) error {
	return &PackageNotFoundError{
		kind:   kind,
		goos:   goos,
		goarch: goarch,
	}
}

func (e PackageNotFoundError) Error() string {
	return fmt.Sprintf("Package not found [%s,%s,%s]", e.goos, e.goarch, e.kind)
}

// VersionNotFoundError 版本不存在错误
type VersionNotFoundError struct {
	version string
	goos    string
	goarch  string
}

func IsVersionNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*VersionNotFoundError)
	return ok
}

func NewVersionNotFoundError(version, goos, goarch string) error {
	return &VersionNotFoundError{
		version: version,
		goos:    goos,
		goarch:  goarch,
	}
}

func (e VersionNotFoundError) Error() string {
	return fmt.Sprintf("Version not found %q [%s,%s]", e.version, e.goos, e.goarch)
}

func (e VersionNotFoundError) Version() string {
	return e.version
}

// MalformedVersionError 版本号格式错误
type MalformedVersionError struct {
	err     error
	version string
}

func NewMalformedVersionError(version string, err error) error {
	return &MalformedVersionError{
		err:     err,
		version: version,
	}
}

func (e MalformedVersionError) Error() string {
	return fmt.Sprintf("Malformed version string %q", e.version)
}

func (e MalformedVersionError) Err() error {
	return e.err
}

func (e MalformedVersionError) Version() string {
	return e.version
}

// URLUnreachableError URL不可达错误
type URLUnreachableError struct {
	err error
	url string
}

// NewURLUnreachableError 返回URL不可达错误实例
func NewURLUnreachableError(url string, err error) error {
	return &URLUnreachableError{
		err: err,
		url: url,
	}
}

func (e URLUnreachableError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("URL %q is unreachable", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

func (e URLUnreachableError) Err() error {
	return e.err
}

func (e URLUnreachableError) URL() string {
	return e.url
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

// Error 返回错误字符串
func (e DownloadError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Resource(%s) download failed", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

// Err 返回错误对象
func (e DownloadError) Err() error {
	return e.err
}

// URL 返回资源URL
func (e DownloadError) URL() string {
	return e.url
}
