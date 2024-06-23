package aliyun

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/voidint/g/collector"
	"github.com/voidint/g/pkg/errs"
	httppkg "github.com/voidint/g/pkg/http"
	"github.com/voidint/g/version"
)

// var _ collector.Collector = (*Collector)(nil)

const (
	// DownloadPageURL 阿里云镜像站点网址
	DownloadPageDomain = "mirrors.aliyun.com"
	DownloadPageURL    = "https://" + DownloadPageDomain + "/golang/"
)

func init() {
	collector.Register(DownloadPageDomain, NewCollector) // TODO 一个采集器对应多个域名
}

// Collector 阿里云镜像站点版本采集器
type Collector struct {
	url  string
	pURL *url.URL
	doc  *goquery.Document
}

// NewCollector 返回采集器实例
func NewCollector() (collector.Collector, error) {
	pURL, err := url.Parse(DownloadPageURL)
	if err != nil {
		return nil, err
	}

	c := Collector{
		url:  DownloadPageURL,
		pURL: pURL,
	}
	if err = c.loadDocument(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Collector) loadDocument() (err error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return errs.NewURLUnreachableError(c.url, err)
	}
	defer resp.Body.Close()

	if !httppkg.IsSuccess(resp.StatusCode) {
		return errs.NewURLUnreachableError(c.url, fmt.Errorf("%d", resp.StatusCode))
	}

	c.doc, err = goquery.NewDocumentFromReader(resp.Body)
	return err
}

// StableVersions 返回所有稳定版本
func (c *Collector) StableVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

// UnstableVersions 返回所有非稳定版本
func (c *Collector) UnstableVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

// ArchivedVersions 返回已归档版本
func (c *Collector) ArchivedVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil
}

// AllVersions 返回所有已知版本
func (c *Collector) AllVersions() (vers []*version.Version, err error) {
	items := c.findGoFileItems(c.doc.Find(".table"))
	if len(items) == 0 {
		return make([]*version.Version, 0), nil
	}
	if vers, err = convert2Versions(items); err != nil {
		return nil, err
	}
	return vers, nil
}

func (c *Collector) findGoFileItems(table *goquery.Selection) (items []*goFileItem) {
	trs := table.Find("tbody").Find("tr")
	items = make([]*goFileItem, 0, trs.Length())

	trs.Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		if !strings.HasPrefix(href, "go") {
			return
		}

		items = append(items, &goFileItem{
			FileName: td.Eq(0).Find("a").Text(),
			URL:      c.url + href,
			Size:     td.Eq(1).Text(),
		})
	})
	return items
}
