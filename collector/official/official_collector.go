package official

import (
	"fmt"
	"net/http"
	stdurl "net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/voidint/g/pkg/errs"
	httppkg "github.com/voidint/g/pkg/http"
	"github.com/voidint/g/version"
)

const (
	// Name Collector name
	Name = "official"
)

// Collector 官方站点版本采集器
type Collector struct {
	url  string
	pURL *stdurl.URL
	doc  *goquery.Document
}

// NewCollector 返回采集器实例
func NewCollector(downloadPageURL string) (*Collector, error) {
	if downloadPageURL == "" {
		return nil, errs.ErrEmptyURL
	}

	pURL, err := stdurl.Parse(downloadPageURL)
	if err != nil {
		return nil, err
	}

	c := Collector{
		url:  downloadPageURL,
		pURL: pURL,
	}
	if err = c.loadDocument(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Name Collector name
func (c *Collector) Name() string {
	return Name
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
			Kind:      version.PackageKind(td.Eq(1).Text()),
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

// StableVersions Return all stable versions
func (c *Collector) StableVersions() (items []*version.Version, err error) {
	var divs *goquery.Selection
	if c.hasUnstableVersions() {
		divs = c.doc.Find("#stable").NextUntil("#unstable")
	} else {
		divs = c.doc.Find("#stable").NextUntil("#archive")
	}

	divs.EachWithBreak(func(i int, div *goquery.Selection) bool {
		vname, ok := div.Attr("id")
		if !ok {
			return true
		}

		var v *version.Version
		if v, err = version.New(
			strings.TrimPrefix(vname, "go"),
			version.WithPackages(c.findPackages(div.Find("table").First())),
		); err != nil {
			return false
		}

		items = append(items, v)
		return true
	})

	if err != nil {
		return nil, err
	}
	sort.Sort(version.Collection(items))
	return items, nil
}

// UnstableVersions Return all stable versions
func (c *Collector) UnstableVersions() (items []*version.Version, err error) {
	c.doc.Find("#unstable").NextUntil("#archive").EachWithBreak(func(i int, div *goquery.Selection) bool {
		vname, ok := div.Attr("id")
		if !ok {
			return true
		}

		var v *version.Version
		if v, err = version.New(
			strings.TrimPrefix(vname, "go"),
			version.WithPackages(c.findPackages(div.Find("table").First())),
		); err != nil {
			return false
		}

		items = append(items, v)
		return true
	})

	if err != nil {
		return nil, err
	}
	sort.Sort(version.Collection(items))
	return items, nil
}

// ArchivedVersions Return all archived versions
func (c *Collector) ArchivedVersions() (items []*version.Version, err error) {
	c.doc.Find("#archive").Find("div.toggle").EachWithBreak(func(i int, div *goquery.Selection) bool {
		vname, ok := div.Attr("id")
		if !ok {
			return true
		}

		var v *version.Version
		if v, err = version.New(
			strings.TrimPrefix(vname, "go"),
			version.WithPackages(c.findPackages(div.Find("table").First())),
		); err != nil {
			return false
		}

		items = append(items, v)
		return true
	})

	if err != nil {
		return nil, err
	}
	sort.Sort(version.Collection(items))
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
	sort.Sort(version.Collection(items))
	return items, nil
}
