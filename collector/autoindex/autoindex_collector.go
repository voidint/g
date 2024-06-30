package autoindex

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/voidint/g/collector/internal"
	"github.com/voidint/g/pkg/errs"
	httppkg "github.com/voidint/g/pkg/http"
	"github.com/voidint/g/version"
)

const (
	// Name Collector name
	Name = "autoindex"
)

// Collector Nginx autoindex collector.
type Collector struct {
	url  string
	pURL *url.URL
	doc  *goquery.Document
}

// NewCollector Get the collector instance
func NewCollector(downloadPageURL string) (*Collector, error) {
	if downloadPageURL == "" {
		return nil, errs.ErrEmptyURL
	}

	pURL, err := url.Parse(downloadPageURL)
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

// StableVersions Return all stable versions
func (c *Collector) StableVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil // Unable to determine which versions are stable
}

// UnstableVersions Return all stable versions
func (c *Collector) UnstableVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil // Unable to determine which versions are unstable
}

// ArchivedVersions Return all archived versions
func (c *Collector) ArchivedVersions() (items []*version.Version, err error) {
	return make([]*version.Version, 0), nil // Unable to determine which versions are archived
}

// AllVersions Return all versions
func (c *Collector) AllVersions() (vers []*version.Version, err error) {
	items := c.findGoFileItems()
	if len(items) == 0 {
		return make([]*version.Version, 0), nil
	}
	if vers, err = internal.Convert2Versions(items); err != nil {
		return nil, err
	}
	return vers, nil
}

func (c *Collector) findGoFileItems() (items []*internal.GoFileItem) {
	anchors := c.doc.Find("pre").Find("a")
	items = make([]*internal.GoFileItem, 0, anchors.Length())

	anchors.Each(func(j int, anchor *goquery.Selection) {
		href := anchor.AttrOr("href", "")
		if !strings.HasPrefix(href, "go") || strings.HasSuffix(href, "/") {
			return
		}

		var size string
		if fields := strings.Fields(strings.TrimSpace(anchor.Nodes[0].NextSibling.Data)); len(fields) > 0 {
			size = fields[len(fields)-1]
		}

		items = append(items, &internal.GoFileItem{
			FileName: anchor.Text(),
			URL:      c.url + href,
			Size:     size,
		})
	})

	return items
}
