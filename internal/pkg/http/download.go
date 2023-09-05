package http

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"

	"github.com/voidint/g/internal/pkg/errs"
)

// Download 下载资源并另存为
func Download(srcURL string, filename string, flag int, perm fs.FileMode, withProgress bool) (size int64, err error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		return 0, errs.NewDownloadError(srcURL, err)
	}
	defer resp.Body.Close()

	if !IsSuccess(resp.StatusCode) {
		return 0, errs.NewURLUnreachableError(srcURL, fmt.Errorf("%d", resp.StatusCode))
	}

	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return 0, errs.NewDownloadError(srcURL, err)
	}
	defer f.Close()

	var dst io.Writer
	if withProgress {
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
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

// DownloadAsBytes 返回下载资源的原始字节切片
func DownloadAsBytes(srcURL string) (data []byte, err error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// IsSuccess 返回 http 请求是否成功
func IsSuccess(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
