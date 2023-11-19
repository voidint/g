package collector

import (
	"strings"

	"github.com/voidint/g/collector/aliyun"
	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/version"
)

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
	for i := range urls {
		urls[i] = strings.TrimSpace(urls[i])

		if urls[i] != "" && (strings.HasPrefix(aliyun.DownloadPageURL, urls[i]) || strings.HasPrefix(urls[i], aliyun.DownloadPageURL)) {
			if c, err = aliyun.NewCollector(); err == nil {
				return c, nil
			}
		}

		if c, err = official.NewCollector(urls[i]); err == nil {
			return c, nil
		}
	}
	return c, err
}
