package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Masterminds/semver"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_inuse(t *testing.T) {
	Convey("查询当前使用中的go版本", t, func() {
		rootDir := filepath.Join(os.TempDir(), fmt.Sprintf(".g_%d", time.Now().Unix()))
		goroot = filepath.Join(rootDir, "go")
		versionsDir = filepath.Join(rootDir, "versions")
		vDir := filepath.Join(versionsDir, "1.12.6")

		os.MkdirAll(versionsDir, 0755)
		os.MkdirAll(vDir, 0755)
		defer os.RemoveAll(rootDir)

		So(os.Symlink(vDir, goroot), ShouldBeNil)
		So(inuse(goroot), ShouldEqual, "1.12.6")
	})
}

func Test_render(t *testing.T) {
	Convey("渲染go版本列表", t, func() {
		var buf strings.Builder
		v1, _ := semver.NewVersion("1.11.11")
		v2, _ := semver.NewVersion("1.7.0")
		v3, _ := semver.NewVersion("1.8.1")
		items := []*semver.Version{v1, v2, v3}

		render("1.8.1", items, &buf)
		So(buf.String(), ShouldEqual, "  1.7\n* 1.8.1\n  1.11.11\n")
	})
}
