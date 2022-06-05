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

// NewCollector 返回采集器实例
func NewCollector(url string) (Collector, error) {
	if strings.HasPrefix(aliyun.DownloadPageURL, url) {
		return aliyun.NewCollector()
	}
	return official.NewCollector(url)
}
