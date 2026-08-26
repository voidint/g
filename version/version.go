// Copyright (c) 2019 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/voidint/g/pkg/checksum"
	"github.com/voidint/g/pkg/errs"
	httppkg "github.com/voidint/g/pkg/http"
)

// Semantify converts Go version strings to semantic version format with adaptation.
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

// Version represents a Go language distribution version.
type Version struct {
	name string // Original version name (e.g. '1.12.4'), may differ from semver format
	sv   *semver.Version
	pkgs []*Package
}

// WithPackages configures available distribution packages for the version.
func WithPackages(pkgs []*Package) func(v *Version) {
	return func(v *Version) {
		v.pkgs = pkgs
	}
}

// New creates a Version instance with semantic version validation.
func New(name string, opts ...func(v *Version)) (*Version, error) {
	sv, err := Semantify(name)
	if err != nil {
		return nil, err
	}

	v := Version{
		name: name,
		sv:   sv,
	}

	for _, setter := range opts {
		if setter == nil {
			continue
		}
		setter(&v)
	}

	return &v, nil
}

func MustNew(name string, opts ...func(v *Version)) *Version {
	v, err := New(name, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// Name returns original version string (e.g. 'go1.21.4').
func (v *Version) Name() string {
	return v.name
}

// Major returns the major number of the semantic version.
func (v *Version) Major() uint64 {
	return v.sv.Major()
}

// Minor returns the minor number of the semantic version.
func (v *Version) Minor() uint64 {
	return v.sv.Minor()
}

// Packages returns all distribution packages for different OS/ARCH combinations.
func (v *Version) Packages() []Package {
	items := make([]Package, 0, len(v.pkgs))
	for _, pkg := range v.pkgs {
		items = append(items, *pkg)
	}
	return items
}

// MatchConstraint checks if version satisfies semantic version constraints.
func (v *Version) MatchConstraint(c *semver.Constraints) bool {
	return c.Check(v.sv)
}

func (v *Version) match(goos, goarch string) bool {
	for _, pkg := range v.pkgs {
		if strings.Contains(pkg.FileName, goos) && strings.Contains(pkg.FileName, goarch) { // TODO: Improve architecture matching logic
			return true
		}
	}
	return false
}

// FindPackages discovers packages matching specific OS/ARCH and package type.
func (v *Version) FindPackages(kind PackageKind, goos, goarch string) (pkgs []Package, err error) {
	prefix := fmt.Sprintf("go%s.%s-%s", v.name, goos, goarch)
	for i := range v.pkgs {
		if v.pkgs[i] == nil || !strings.EqualFold(string(v.pkgs[i].Kind), string(kind)) || !strings.HasPrefix(v.pkgs[i].FileName, prefix) {
			continue
		}
		pkgs = append(pkgs, *v.pkgs[i])
	}
	if len(pkgs) == 0 {
		return nil, errs.NewPackageNotFoundError(string(kind), goos, goarch)
	}
	return pkgs, nil
}

// Package describes a Go distribution file metadata.
type Package struct {
	FileName    string      `json:"filename"`
	URL         string      `json:"url"`
	Kind        PackageKind `json:"kind"`
	OS          string      `json:"os"`
	Arch        string      `json:"arch"`
	Size        string      `json:"size"`
	Checksum    string      `json:"checksum"`
	ChecksumURL string      `json:"-"`
	Algorithm   string      `json:"algorithm"` // checksum algorithm
}

// PackageKind indicates distribution package format type.
type PackageKind string

const (
	// SourceKind indicates source code package
	SourceKind PackageKind = "Source"
	// ArchiveKind indicates compressed archive package
	ArchiveKind PackageKind = "Archive"
	// InstallerKind indicates executable installer
	InstallerKind PackageKind = "Installer"
)

// DownloadWithProgress fetches package with real-time download metrics.
func (pkg *Package) DownloadWithProgress(dst string) (size int64, err error) {
	return httppkg.Download(pkg.URL, dst, os.O_CREATE|os.O_WRONLY, 0644, true)
}

// VerifyChecksum validates downloaded file against cryptographic hash.
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
