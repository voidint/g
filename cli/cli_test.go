package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func Test_ghome(t *testing.T) {
	t.Run("查询ghome路径", func(t *testing.T) {
		home, err := os.UserHomeDir()
		assert.Nil(t, err)
		assert.Equal(t, filepath.Join(home, ".g"), ghome())
	})
}

func Test_inuse(t *testing.T) {
	t.Run("查询当前使用中的go版本", func(t *testing.T) {
		rootDir := filepath.Join(os.TempDir(), fmt.Sprintf(".g_%d", time.Now().Unix()))
		goroot = filepath.Join(rootDir, "go")
		versionsDir = filepath.Join(rootDir, "versions")
		vDir := filepath.Join(versionsDir, "1.12.6")

		_ = os.MkdirAll(versionsDir, 0755)
		_ = os.MkdirAll(vDir, 0755)
		defer os.RemoveAll(rootDir)

		assert.Nil(t, mkSymlink(vDir, goroot))
		assert.Equal(t, "1.12.6", inuse(goroot))
	})
}

func Test_render(t *testing.T) {
	t.Run("渲染go版本列表", func(t *testing.T) {
		var buf strings.Builder
		v0, _ := semver.NewVersion("1.13-beta1")
		v1, _ := semver.NewVersion("1.11.11")
		v2, _ := semver.NewVersion("1.7.0")
		v3, _ := semver.NewVersion("1.8.1")
		items := []*semver.Version{v0, v1, v2, v3}

		render(map[string]bool{"1.8.1": true}, items, &buf)
		assert.Equal(t, "  1.7\n* 1.8.1\n  1.11.11\n  1.13beta1\n", buf.String())
	})
}

func Test_wrapstring(t *testing.T) {
	t.Run("包装字符串", func(t *testing.T) {
		assert.Equal(t, "[g] Hello world", wrapstring("hello world"))
	})
}

func Test_errstring(t *testing.T) {
	t.Run("返回错误字符串", func(t *testing.T) {
		assert.Equal(t, "", errstring(nil))
		assert.Equal(t, "[g] Hello world", errstring(errors.New("hello world")))
	})
}
