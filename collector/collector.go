package collector

import (
	stdurl "net/url"
	"strings"

	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/version"
)

var collectors = make(map[string]Builder)

func Register(domain string, b Builder) {
	collectors[domain] = b
}

type Builder func() (Collector, error)

// Collector 版本信息采集器
type Collector interface {
	// 返回稳定版本列表
	StableVersions() (items []*version.Version, err error)
	// 返回非稳定版本列表
	UnstableVersions() (items []*version.Version, err error)
	// 返回已归档版本列表
	ArchivedVersions() (items []*version.Version, err error)
	// 返回所有版本列表
	AllVersions() (items []*version.Version, err error)
}

// NewCollector 返回首个可用的采集器实例
func NewCollector(urls ...string) (c Collector, err error) {
	if len(urls) == 0 {
		urls = []string{official.DefaultDownloadPageURL}
	}
	for _, rawUrl := range urls {
		var url *stdurl.URL
		url, err = stdurl.Parse(strings.TrimSpace(rawUrl))
		if err != nil {
			continue
		}

		for domain, c := range collectors {
			if url.Host == strings.ToLower(domain) {
				return c()
			}
		}

		if c, err = official.NewCollector(rawUrl); err == nil {
			return c, nil
		}
	}
	return c, err
}
