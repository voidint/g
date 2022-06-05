package aliyun

import (
	"net/http"
	stdurl "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/voidint/g/pkg/errs"
	"github.com/voidint/g/version"
)

// var _ collector.Collector = (*Collector)(nil)

const (
	// DownloadPageURL 阿里云镜像站点网址
	DownloadPageURL = "https://mirrors.aliyun.com/golang/"
)

// Collector 阿里云镜像站点版本采集器
type Collector struct {
	url  string
	pURL *stdurl.URL
	doc  *goquery.Document
}

// NewCollector 返回采集器实例
func NewCollector() (*Collector, error) {
	pURL, err := stdurl.Parse(DownloadPageURL)
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
	if resp.StatusCode != http.StatusOK {
		return errs.NewURLUnreachableError(c.url, nil)
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
		return make([]*version.Version, 0, 0), nil
	}
	return convert2Versions(items), nil
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
