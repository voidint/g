package ustc

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
	// DownloadPageURL USTC镜像站点网址
	DownloadPageDomain = "mirrors.ustc.edu.cn"
	DownloadPageURL    = "https://" + DownloadPageDomain + "/golang/"
)

func init() {
	collector.Register(DownloadPageDomain, NewCollector) // TODO 一个采集器对应多个域名
}

// Collector USTC镜像站点版本采集器
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
	items := c.findGoFileItems()
	if len(items) == 0 {
		return make([]*version.Version, 0), nil
	}
	if vers, err = convert2Versions(items); err != nil {
		return nil, err
	}
	return vers, nil
}

func (c *Collector) findGoFileItems() (items []*goFileItem) {
	anchors := c.doc.Find("a")

	items = make([]*goFileItem, 0, anchors.Length())

	anchors.Each(func(j int, anchor *goquery.Selection) {
		href := anchor.AttrOr("href", "")
		if !strings.HasPrefix(href, "go") || strings.HasSuffix(href, "/") {
			return
		}

		datas := strings.Split(anchor.Nodes[0].NextSibling.Data, " ")
		size := datas[len(datas)-1]
		size = size[:len(size)-1]

		items = append(items, &goFileItem{
			FileName: anchor.Text(),
			URL:      c.url + href,
			Size:     size,
		})
	})

	return items
}
