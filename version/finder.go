package version

import (
	"runtime"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/voidint/g/pkg/errs"
)

// Finder 版本查找器
type Finder struct {
	goos   string
	goarch string
	items  []*Version
}

// WithFinderGoos 设置查找器所在的目标操作系统，如darwin, freebsd, linux等。
func WithFinderGoos(goos string) func(fdr *Finder) {
	return func(fdr *Finder) {
		fdr.goos = goos
	}
}

// WithFinderGoarch 设置查找器所在的目标硬件架构，如386, amd64, arm, s390x等。
func WithFinderGoarch(goarch string) func(fdr *Finder) {
	return func(fdr *Finder) {
		fdr.goarch = goarch
	}
}

// NewFinder 返回
func NewFinder(items []*Version, opts ...func(fdr *Finder)) *Finder {
	sort.Sort(Collection(items)) // 升序

	fdr := Finder{
		goos:   runtime.GOOS,
		goarch: runtime.GOARCH,
		items:  items,
	}

	if opts != nil {
		for _, setter := range opts {
			setter(&fdr)
		}
	}

	return &fdr
}

// Find 返回满足条件的语义化版本号。版本格式：主版本号.次版本号.修订号。
// vname 支持以下几类版本标识：
// 1、具体版本号：如'1.21.4'
// 2、最新版本：latest
// 3、通配符：如'1.21.x'、'1.x'、'1.18.*'等
// 4、兼容某个主版本号：如'^1'、'^1.18'、'^1.18.10'等，在主版本号保持一致的前提下，次版本号和修订号均保持最新。
// 5、匹配某个主次版本号：如'~1.18'，在主次版本号保持一致的前提下，修订号保持最新。
// 6、大于某个版本：如'>1.18'，大于该版本的前提下，匹配最大的版本号。
// 7、小于某个版本：如'<1.16'，小于该版本的前提下，匹配最大的版本号。
// 8、版本区间：如'1.18 - 1.20'，匹配该区间范围内的最大版本。
func (fdr *Finder) Find(vname string) (*Version, error) {
	if vname == latest {
		return fdr.findLatest()
	}

	for i := len(fdr.items) - 1; i > 0; i-- {
		if fdr.items[i].name == vname && fdr.items[i].match(fdr.goos, fdr.goarch) {
			return fdr.items[i], nil
		}
	}

	cs, err := semver.NewConstraint(vname)
	if err != nil {
		return nil, errs.NewVersionNotFoundError(vname, fdr.goos, fdr.goarch)
	}

	for i := len(fdr.items) - 1; i > 0; i-- { // 优先匹配高版本
		if fdr.items[i].match(fdr.goos, fdr.goarch) && cs.Check(fdr.items[i].sv) {
			return fdr.items[i], nil
		}
	}
	return nil, errs.NewVersionNotFoundError(vname, fdr.goos, fdr.goarch)
}

const latest = "latest"

func (fdr *Finder) findLatest() (*Version, error) {
	for i := len(fdr.items) - 1; i > 0; i-- {
		if fdr.items[i].match(fdr.goos, fdr.goarch) {
			return fdr.items[i], nil
		}
	}
	return nil, errs.NewVersionNotFoundError(latest, fdr.goos, fdr.goarch)
}
