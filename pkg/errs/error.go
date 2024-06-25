package errs

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrUnsupportedChecksumAlgorithm Unsupported checksum algorithm
	ErrUnsupportedChecksumAlgorithm = errors.New("unsupported checksum algorithm")
	// ErrChecksumNotMatched File checksum does not match the computed checksum
	ErrChecksumNotMatched = errors.New("file checksum does not match the computed checksum")
	// ErrChecksumFileNotFound Checksum file not found
	ErrChecksumFileNotFound = errors.New("checksum file not found")
	// ErrAssetNotFound Asset not found
	ErrAssetNotFound = errors.New("asset not found")
	// ErrCollectorNotFound Collector not found
	ErrCollectorNotFound = errors.New("collector not found")
	// ErrEmptyURL URL is empty
	ErrEmptyURL = errors.New("empty url")
)

// PackageNotFoundError 软件包不存在错误
type PackageNotFoundError struct {
	kind   string
	goos   string
	goarch string
}

// IsPackageNotFound 若是软件包不存在错误，则返回true；反之，返回false。
func IsPackageNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*PackageNotFoundError)
	return ok
}

// NewPackageNotFoundError 返回软件包不存在错误实例
func NewPackageNotFoundError(kind, goos, goarch string) error {
	return &PackageNotFoundError{
		kind:   kind,
		goos:   goos,
		goarch: goarch,
	}
}

// Error 返回错误详情
func (e PackageNotFoundError) Error() string {
	return fmt.Sprintf("package not found [%s,%s,%s]", e.goos, e.goarch, e.kind)
}

// VersionNotFoundError 版本不存在错误
type VersionNotFoundError struct {
	version string
	goos    string
	goarch  string
}

// IsVersionNotFound 若是版本不存在错误，返回true；反之，返回false。
func IsVersionNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*VersionNotFoundError)
	return ok
}

// NewVersionNotFoundError 返回版本不存在错误实例
func NewVersionNotFoundError(version, goos, goarch string) error {
	return &VersionNotFoundError{
		version: version,
		goos:    goos,
		goarch:  goarch,
	}
}

// Error 返回错误详情
func (e VersionNotFoundError) Error() string {
	return fmt.Sprintf("version not found %q [%s,%s]", e.version, e.goos, e.goarch)
}

// Version 返回版本号
func (e VersionNotFoundError) Version() string {
	return e.version
}

// MalformedVersionError 版本号格式错误
type MalformedVersionError struct {
	err     error
	version string
}

// IsMalformedVersion 若是版本号格式错误，返回true；反之，返回false。
func IsMalformedVersion(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*MalformedVersionError)
	return ok
}

// NewMalformedVersionError 返回版本号格式错误实例
func NewMalformedVersionError(version string, err error) error {
	return &MalformedVersionError{
		err:     err,
		version: version,
	}
}

// Error 返回错误详情
func (e MalformedVersionError) Error() string {
	return fmt.Sprintf("malformed version string %q", e.version)
}

// Unwrap 返回源错误
func (e MalformedVersionError) Unwrap() error {
	return e.err
}

// Version 返回版本号
func (e MalformedVersionError) Version() string {
	return e.version
}

// URLUnreachableError URL不可达错误
type URLUnreachableError struct {
	err error
	url string
}

// IsURLUnreachable 若是URL不可达错误，返回true；反之，返回false。
func IsURLUnreachable(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*URLUnreachableError)
	return ok
}

// NewURLUnreachableError 返回URL不可达错误实例
func NewURLUnreachableError(url string, err error) error {
	return &URLUnreachableError{
		err: err,
		url: url,
	}
}

// Error 返回错误详情
func (e URLUnreachableError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("URL %q is unreachable", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

// Unwrap 返回源错误
func (e URLUnreachableError) Unwrap() error {
	return e.err
}

// URL 返回URL地址
func (e URLUnreachableError) URL() string {
	return e.url
}

// DownloadError 下载失败错误
type DownloadError struct {
	url string
	err error
}

// IsDownload 若是下载失败错误，返回true；反之，返回false。
func IsDownload(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*DownloadError)
	return ok
}

// NewDownloadError 返回下载失败错误实例
func NewDownloadError(url string, err error) error {
	return &DownloadError{
		url: url,
		err: err,
	}
}

// Error 返回错误详情
func (e DownloadError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("resource(%s) download failed", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

// Unwrap 返回错误对象
func (e DownloadError) Unwrap() error {
	return e.err
}

// URL 返回资源URL
func (e DownloadError) URL() string {
	return e.url
}
