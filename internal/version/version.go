package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/voidint/g/internal/pkg/checksum"
	"github.com/voidint/g/internal/pkg/errs"
	httppkg "github.com/voidint/g/internal/pkg/http"
)

// Semantify go 版本号并未完全遵循语义化版本号标准，该函数进行了一定的适配，返回一个语义化版本。
func Semantify(vname string) (*semver.Version, error) {
	var idx int
	if strings.Contains(vname, "alpha") {
		idx = strings.Index(vname, "alpha")

	} else if strings.Contains(vname, "beta") {
		idx = strings.Index(vname, "beta")

	} else if strings.Contains(vname, "rc") {
		idx = strings.Index(vname, "rc")
	}
	if idx > 0 {
		vname = vname[:idx] + "-" + vname[idx:]
	}

	sv, err := semver.NewVersion(vname)
	if err != nil {
		return nil, errs.NewMalformedVersionError(vname, err)
	}
	return sv, nil
}

// FindVersion 返回指定名称的版本
func FindVersion(all []*Version, name string) (*Version, error) {
	for i := range all {
		if all[i].Name == name {
			return all[i], nil
		}
	}
	return nil, errs.NewVersionNotFoundError(name)
}

// Version go版本
type Version struct {
	Name            string // 版本名，如'1.12.4'
	SemanticVersion *semver.Version
	Packages        []*Package
}

func New(name string, pkgs []*Package) (*Version, error) {
	sv, err := Semantify(name)
	if err != nil {
		return nil, err
	}

	return &Version{
		Name:            name,
		SemanticVersion: sv,
		Packages:        pkgs,
	}, nil
}

func MustNew(name string, pkgs []*Package) *Version {
	v, err := New(name, pkgs)
	if err != nil {
		panic(err)
	}
	return v
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
	FileName    string
	URL         string
	Kind        string
	OS          string
	Arch        string
	Size        string
	Checksum    string
	ChecksumURL string
	Algorithm   string // checksum algorithm
}

const (
	// SourceKind go安装包种类-源码
	SourceKind = "Source"
	// ArchiveKind go安装包种类-压缩文件
	ArchiveKind = "Archive"
	// InstallerKind go安装包种类-可安装程序
	InstallerKind = "Installer"
)

// DownloadWithProgress 下载版本另存为指定文件且显示下载进度
func (pkg *Package) DownloadWithProgress(dst string) (size int64, err error) {
	return httppkg.Download(pkg.URL, dst, os.O_CREATE|os.O_WRONLY, 0644, true)
}

// VerifyChecksum 验证目标文件的校验和与当前安装包的校验和是否一致
func (pkg *Package) VerifyChecksum(filename string) (err error) {
	if pkg.Checksum == "" && pkg.ChecksumURL != "" {
		data, err := httppkg.DownloadAsBytes(pkg.ChecksumURL)
		if err != nil {
			return err
		}
		pkg.Checksum = string(data)
	}
	var algo checksum.Algorithm
	switch pkg.Algorithm {
	case string(checksum.SHA256):
		algo = checksum.SHA256
	case string(checksum.SHA1):
		algo = checksum.SHA1
	default:
		return errs.ErrUnsupportedChecksumAlgorithm
	}
	return checksum.VerifyFile(algo, pkg.Checksum, filename)
}
