package aliyun

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func getCollector() (*Collector, error) {
	b, err := ioutil.ReadFile("./testdata/golang_dl.html")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return &Collector{
		url: DownloadPageURL,
		doc: doc,
	}, nil
}

func Test_findGoFileItems(t *testing.T) {
	c, err := getCollector()
	assert.Nil(t, err)
	assert.NotNil(t, c)

	t.Run("", func(t *testing.T) {
		items := c.findGoFileItems(c.doc.Find(".table"))
		assert.NotEmpty(t, items)
	})
}
