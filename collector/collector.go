package collector

import (
	"strings"

	"github.com/voidint/g/collector/autoindex"
	"github.com/voidint/g/collector/fancyindex"
	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/pkg/errs"
	"github.com/voidint/g/version"
)

// official collector
const (
	// OriginalOfficialDownloadPageURL Golang official site download page URL
	OriginalOfficialDownloadPageURL = "https://golang.org/dl/"
	// OfficialDownloadPageURL Golang official site download page URL
	OfficialDownloadPageURL = "https://go.dev/dl/"
	// CNDownloadPageURL China mirror site download page URL
	CNDownloadPageURL = "https://golang.google.cn/dl/"
)

// Nginx fancyindex collector
const (
	// AliYunDownloadPageURL Alibaba cloud mirror site URL
	AliYunDownloadPageURL = "https://mirrors.aliyun.com/golang/"
	// HUSTDownloadPageURL Huazhong University of Science and Technology mirror site URL
	HUSTDownloadPageURL = "https://mirrors.hust.edu.cn/golang/"
	// NJUDownloadPageURL Nanjing University mirror site URL
	NJUDownloadPageURL = "https://mirrors.nju.edu.cn/golang/"
)

// Nginx autoindex collector
const (
	// USTCDownloadPageURL University of Science and Technology of China mirror site URL
	USTCDownloadPageURL = "https://mirrors.ustc.edu.cn/golang/"
)

// Collector Version information collector
type Collector interface {
	// Name Collector name
	Name() string
	// StableVersions Return all stable versions
	StableVersions() (items []*version.Version, err error)
	// UnstableVersions Return all stable versions
	UnstableVersions() (items []*version.Version, err error)
	// ArchivedVersions Return all archived versions
	ArchivedVersions() (items []*version.Version, err error)
	// AllVersions Return all versions
	AllVersions() (items []*version.Version, err error)
}

// NewCollector Returns the first available collector instance
// official|https://go.dev/dl/,fancyindex|https://mirrors.aliyun.com/golang/,autoindex|https://mirrors.ustc.edu.cn/golang/
func NewCollector(urls ...string) (c Collector, err error) {
	if size := len(urls); size == 0 || (size == 1 && urls[0] == "") {
		urls = []string{OfficialDownloadPageURL}
	}

	for i := range urls {
		urls[i] = strings.TrimSpace(urls[i])

		if !strings.HasSuffix(urls[i], "/") {
			urls[i] = urls[i] + "/"
		}

		idx := strings.Index(urls[i], "|")

		if idx > 0 && idx < len(urls[i])-1 {
			downloadPageURL := strings.TrimSpace(urls[i][idx+1:])

			switch collectorName := strings.TrimSpace(urls[i][:idx]); collectorName {
			case official.Name:
				return official.NewCollector(downloadPageURL)

			case fancyindex.Name:
				return fancyindex.NewCollector(downloadPageURL)

			case autoindex.Name:
				return autoindex.NewCollector(downloadPageURL)

			default:
				continue
			}
		}

		switch urls[i] {
		case OfficialDownloadPageURL, OriginalOfficialDownloadPageURL, CNDownloadPageURL:
			return official.NewCollector(urls[i])

		case AliYunDownloadPageURL, HUSTDownloadPageURL, NJUDownloadPageURL:
			return fancyindex.NewCollector(urls[i])

		case USTCDownloadPageURL:
			return autoindex.NewCollector(urls[i])

		default:
			continue
		}
	}
	return nil, errs.ErrCollectorNotFound
}
