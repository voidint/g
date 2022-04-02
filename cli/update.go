package cli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/urfave/cli"
	"github.com/voidint/g/build"
	"github.com/voidint/g/github"
)

func update(ctx *cli.Context) (err error) {
	up := github.NewReleaseUpdater(findAsset)

	// 检查更新
	latest, yes, err := up.CheckForUpdates(semver.MustParse(build.ShortVersion), "voidint", "g")
	if err != nil {
		return cli.NewExitError(errstring(err), 1)
	}
	if !yes {
		fmt.Printf("You are up to date! g v%s is the latest version.\n", build.ShortVersion)
		return nil
	}
	fmt.Printf("A new version of g(%s) is available\n", latest.TagName)

	// 应用更新
	if err = up.Apply(latest); err != nil {
		return cli.NewExitError(errstring(err), 1)
	}
	fmt.Println("Update completed")
	return nil
}

func findAsset(items []github.Asset) (idx int) {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	suffix := fmt.Sprintf("%s-%s.%s", runtime.GOOS, runtime.GOARCH, ext)
	for i := range items {
		if strings.HasSuffix(items[i].BrowserDownloadURL, suffix) {
			return i
		}
	}
	return -1
}
