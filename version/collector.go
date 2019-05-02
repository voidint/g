package version

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Collector 采集器
type Collector struct {
	url string
	doc *goquery.Document
}

// NewCollector 返回采集器实例
func NewCollector(url string) (*Collector, error) {
	c := Collector{
		url: url,
	}
	if err := c.loadDocument(); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Collector) loadDocument() (err error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return NewURLUnreachableError(c.url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return NewURLUnreachableError(c.url, nil)
	}
	c.doc, err = goquery.NewDocumentFromReader(resp.Body)
	return err
}

func (c *Collector) findPackages(table *goquery.Selection) (pkgs []*Package) {
	table.Find("tr").Not(".first").Each(func(j int, tr *goquery.Selection) {
		td := tr.Find("td")
		pkgs = append(pkgs, &Package{
			FileName: td.Eq(0).Find("a").Text(),
			URL:      td.Eq(0).Find("a").AttrOr("href", ""),
			Kind:     td.Eq(1).Text(),
			OS:       td.Eq(2).Text(),
			Arch:     td.Eq(3).Text(),
			Size:     td.Eq(4).Text(),
			Checksum: td.Eq(5).Text(),
		})
	})
	return pkgs
}

// StableVersions 返回所有稳定版本
func (c *Collector) StableVersions() (items []*Version, err error) {
	c.doc.Find("#stable").NextUntil("#archive").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		items = append(items, &Version{
			Name:     strings.TrimPrefix(vname, "go"),
			Packages: c.findPackages(div.Find("table").First()),
		})
	})
	return items, nil
}

// ArchivedVersions 返回已归档版本
func (c *Collector) ArchivedVersions() (items []*Version, err error) {
	c.doc.Find("#archive").Find("div.toggle").Each(func(i int, div *goquery.Selection) {
		vname, ok := div.Attr("id")
		if !ok {
			return
		}
		items = append(items, &Version{
			Name:     strings.TrimPrefix(vname, "go"),
			Packages: c.findPackages(div.Find("table").First()),
		})
	})
	return items, nil
}

// AllVersions 返回所有已知版本
func (c *Collector) AllVersions() (items []*Version, err error) {
	items, err = c.StableVersions()
	if err != nil {
		return nil, err
	}
	archives, err := c.ArchivedVersions()
	if err != nil {
		return nil, err
	}
	items = append(items, archives...)
	return items, nil
}
