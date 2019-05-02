package version

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// DefaultURL 提供go版本信息的默认网址
	DefaultURL = "https://golang.org/dl/"
)

var (
	// ErrVersionNotFound 版本不存在
	ErrVersionNotFound = errors.New("version not found")
	// ErrPackageNotFound 版本包不存在
	ErrPackageNotFound = errors.New("package not found")
)

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

func (e *URLUnreachableError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("URL %q is unreachable", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

// Package go版本安装包
type Package struct {
	FileName string
	URL      string
	Kind     string
	OS       string
	Arch     string
	Size     string
	Checksum string
}

// Download 下载版本另存为指定文件并校验sha256哈希值
func (pkg *Package) Download(filename string) error {
	// TODO 待实现
	return nil
}

// Version go版本
type Version struct {
	Name     string // 版本名，如'1.12.4'
	Packages []*Package
}

// Download 下载版本另存为指定文件并校验sha256哈希值
func (v *Version) Download(os, arch string, filename string) error {
	// TODO 待实现
	return nil
}

// FindPackage 返回指定操作系统和硬件架构的版本包
func FindPackage(all []*Package, os, arch string) (*Package, error) {
	for i := range all {
		if all[i] != nil && strings.EqualFold(all[i].OS, os) && strings.EqualFold(all[i].Arch, arch) {
			return all[i], nil
		}
	}
	return nil, ErrPackageNotFound
}
