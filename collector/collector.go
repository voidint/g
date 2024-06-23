package collector

import (
	"net/url"
	"strings"
	"sync"

	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/version"
)

type Builder func() (Collector, error)

var (
	mu         sync.RWMutex
	collectors = make(map[string]Builder)
)

func Register(domain string, builder Builder) {
	mu.Lock()
	defer mu.Unlock()

	if builder == nil {
		panic("register builder is nil")
	}

	if _, dup := collectors[domain]; dup {
		panic("register called twice for builder")
	}

	collectors[domain] = builder
}

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
	mu.RLock()
	defer mu.RUnlock()

	if len(urls) == 0 {
		urls = []string{official.DefaultDownloadPageURL}
	}
	for i := range urls {
		pURL, err := url.Parse(strings.TrimSpace(urls[i]))
		if err != nil {
			continue
		}

		for domain, c := range collectors {
			if pURL.Host == strings.ToLower(domain) {
				return c()
			}
		}

		if c, err = official.NewCollector(urls[i]); err == nil {
			return c, nil
		}
	}
	return c, err
}
