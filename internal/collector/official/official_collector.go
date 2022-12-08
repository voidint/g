package official

import (
	"fmt"
	"net/http"
	stdurl "net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/voidint/g/internal/pkg/errs"
	"github.com/voidint/g/internal/version"
)

// var _ collector.Collector = (*Collector)(nil)

const (
	// DefaultDownloadPageURL 官方下载站点网址
	DefaultDownloadPageURL = "https://go.dev/dl/"
)

// Collector 官方站点版本采集器
type Collector struct {
	url  string
	pURL *stdurl.URL
	doc  *goquery.Document
}

// NewCollector 返回采集器实例
func NewCollector(url string) (*Collector, error) {
	if url == "" {
		url = DefaultDownloadPageURL
	}
	pURL, err := stdurl.Parse(url)
	if err != nil {
		return nil, err
	}

	c := Collector{
		url:  url,
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

func (c *Collector) findPackages(table *goquery.Selection) (pkgs []*version.Package) {
	alg := strings.TrimSuffix(table.Find("thead").Find("th").Last().Text(), " Checksum")

	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		href := td.Eq(0).Find("a").AttrOr("href", "")
		if strings.HasPrefix(href, "/") { // relative paths
			href = fmt.Sprintf("%s://%s%s", c.pURL.Scheme, c.pURL.Host, href)
		}
		pkgs = append(pkgs, &version.Package{
			FileName:  td.Eq(0).Find("a").Text(),
			URL:       href,
			Kind:      td.Eq(1).Text(),
			OS:        td.Eq(2).Text(),
			Arch:      td.Eq(3).Text(),
			Size:      td.Eq(4).Text(),
			Checksum:  td.Eq(5).Text(),
			Algorithm: alg,
		})
	})
	return pkgs
}

// hasUnstableVersions 返回是否包含非稳定版本的布尔值
func (c *Collector) hasUnstableVersions() bool {
	return c.doc.Find("#unstable").Length() > 0
}

// StableVersions 返回所有稳定版本
func (c *Collector) StableVersions() (items []*version.Version, err error) {
	var divs *goquery.Selection
	if c.hasUnstableVersions() {
		divs = c.doc.Find("#stable").NextUntil("#unstable")
	} else {
		divs = c.doc.Find("#stable").NextUntil("#archive")
	}

	divs.Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		items = append(items, &version.Version{
			Name:     strings.TrimPrefix(vname, "go"),
			Packages: c.findPackages(div.Find("table").First()),
		})
	})
	return items, nil
}

// UnstableVersions 返回所有非稳定版本
func (c *Collector) UnstableVersions() (items []*version.Version, err error) {
	c.doc.Find("#unstable").NextUntil("#archive").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		items = append(items, &version.Version{
			Name:     strings.TrimPrefix(vname, "go"),
			Packages: c.findPackages(div.Find("table").First()),
		})
	})
	return items, nil
}

// ArchivedVersions 返回已归档版本
func (c *Collector) ArchivedVersions() (items []*version.Version, err error) {
	c.doc.Find("#archive").Find("div.toggle").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		items = append(items, &version.Version{
			Name:     strings.TrimPrefix(vname, "go"),
			Packages: c.findPackages(div.Find("table").First()),
		})
	})
	return items, nil
}

// AllVersions 返回所有已知版本
func (c *Collector) AllVersions() (items []*version.Version, err error) {
	items, err = c.StableVersions()
	if err != nil {
		return nil, err
	}
	archives, err := c.ArchivedVersions()
	if err != nil {
		return nil, err
	}
	items = append(items, archives...)

	unstables, err := c.UnstableVersions()
	if err != nil {
		return nil, err
	}
	items = append(items, unstables...)
	return items, nil
}
