package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLUnreachableError(t *testing.T) {
	t.Run("URL不可达错误", func(t *testing.T) {
		url := "https://github.com/voidint"
		core := errors.New("hello error")

		err := NewURLUnreachableError(url, core)
		assert.NotNil(t, err)

		e, ok := err.(*URLUnreachableError)
		assert.True(t, ok)
		assert.NotNil(t, e)
		assert.Equal(t, url, e.URL())
		assert.Equal(t, core, e.Err())
		assert.Equal(t, fmt.Sprintf("URL %q is unreachable ==> %s", url, core.Error()), e.Error())
	})
}

func TestDownloadError(t *testing.T) {
	t.Run("安装包下载错误", func(t *testing.T) {
		url := "https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz"
		core := errors.New("hello error")

		err := NewDownloadError(url, core)
		assert.NotNil(t, err)
		e, ok := err.(*DownloadError)
		assert.True(t, ok)
		assert.NotNil(t, e)
		assert.Equal(t, url, e.URL())
		assert.Equal(t, core, e.Err())
		assert.Equal(t, fmt.Sprintf("resource(%s) download failed ==> %s", url, core.Error()), e.Error())
	})
}
