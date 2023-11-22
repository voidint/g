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

	for _, setter := range opts {
		if opts != nil {
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
// 4、匹配最新的次版本号（主版本号兼容）：如'^1'、'^1.18'、'^1.18.10'等，在主版本号保持一致的前提下，次版本号和修订号均保持最新。
// 5、匹配某个次版本号的最新修订号：如'~1.18'，在主次版本号保持一致的前提下，修订号保持最新。
// 6、匹配大于目标版本的最新版本：如'>1.18'，大于该版本的前提下，匹配最大的版本号。
// 7、匹配小于目标版本的最新版本：如'<1.16'，小于该版本的前提下，匹配最大的版本号。
// 8、匹配目标版本区间内的最新版本：如'1.18 - 1.20'，匹配该区间范围内的最大版本。
func (fdr *Finder) Find(vname string) (*Version, error) {
	if vname == Latest {
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

// MustFind 返回满足条件的语义化版本号。若发生错误，则抛出panic。
func (fdr *Finder) MustFind(vname string) *Version {
	v, err := fdr.Find(vname)
	if err != nil {
		panic(err)
	}
	return v
}

// Latest 指代当前最新版本
const Latest = "latest"

func (fdr *Finder) findLatest() (*Version, error) {
	for i := len(fdr.items) - 1; i > 0; i-- {
		if fdr.items[i].match(fdr.goos, fdr.goarch) {
			return fdr.items[i], nil
		}
	}
	return nil, errs.NewVersionNotFoundError(Latest, fdr.goos, fdr.goarch)
}
