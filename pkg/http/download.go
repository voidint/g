package http

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// DownloadError 下载失败错误
type DownloadError struct {
	url string
	err error
}

// NewDownloadError 返回下载失败错误实例
func NewDownloadError(url string, err error) error {
	return &DownloadError{
		url: url,
		err: err,
	}
}

// Error 返回错误字符串
func (e DownloadError) Error() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Resource(%s) download failed", e.url))
	if e.err != nil {
		buf.WriteString(" ==> " + e.err.Error())
	}
	return buf.String()
}

// Err 返回错误对象
func (e DownloadError) Err() error {
	return e.err
}

// URL 返回资源URL
func (e DownloadError) URL() string {
	return e.url
}

// Download 下载资源并另存为
func Download(srcURL string, filename string, flag int, perm fs.FileMode, withProgress bool) (size int64, err error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		return 0, NewDownloadError(srcURL, err)
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return 0, NewDownloadError(srcURL, err)
	}
	defer f.Close()

	var dst io.Writer
	if withProgress {
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription("Downloading"),
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionShowBytes(true),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				_, _ = fmt.Fprint(ansi.NewAnsiStdout(), "\n")
			}),
			// progressbar.OptionSpinnerType(35),
			// progressbar.OptionFullWidth(),
		)
		_ = bar.RenderBlank()
		dst = io.MultiWriter(f, bar)

	} else {
		dst = f
	}
	return io.Copy(dst, resp.Body)
}
