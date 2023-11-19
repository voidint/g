package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/build"
	"github.com/voidint/g/pkg/checksum"
	"github.com/voidint/g/pkg/errs"
	httppkg "github.com/voidint/g/pkg/http"
	"github.com/voidint/g/pkg/sdk/github"
)

func selfUpdate(*cli.Context) (err error) {
	up := github.NewReleaseUpdater()

	// 检查更新
	latest, yes, err := up.CheckForUpdates(semver.MustParse(build.ShortVersion), "voidint", "g")
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}
	if !yes {
		fmt.Printf("You are up to date! g v%s is the latest version.\n", build.ShortVersion)
		return nil
	}
	fmt.Printf("A new version of g(%s) is available\n", latest.TagName)

	// 应用更新
	if err = up.Apply(latest, findAsset, findChecksum); err != nil {
		return cli.Exit(errstring(err), 1)
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

func findChecksum(items []github.Asset) (algo checksum.Algorithm, expectedChecksum string, err error) {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	suffix := fmt.Sprintf("%s-%s.%s", runtime.GOOS, runtime.GOARCH, ext)

	var checksumFileURL string
	for i := range items {
		if items[i].Name == "sha256sum.txt" {
			checksumFileURL = items[i].BrowserDownloadURL
			break
		}
	}
	if checksumFileURL == "" {
		return checksum.SHA256, "", errs.ErrChecksumFileNotFound
	}

	resp, err := http.Get(checksumFileURL)
	if err != nil {
		return checksum.SHA256, "", err
	}
	defer resp.Body.Close()

	if !httppkg.IsSuccess(resp.StatusCode) {
		return "", "", errs.NewURLUnreachableError(checksumFileURL, fmt.Errorf("%d", resp.StatusCode))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasSuffix(line, suffix) {
			continue
		}
		return checksum.SHA256, strings.Fields(line)[0], nil
	}
	if err = scanner.Err(); err != nil {
		return checksum.SHA256, "", err
	}
	return checksum.SHA256, "", errs.ErrChecksumFileNotFound
}
