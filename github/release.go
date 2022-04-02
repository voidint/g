package github

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/go-resty/resty/v2"
	"github.com/mholt/archiver/v3"
	"github.com/voidint/g/errs"
	"github.com/voidint/go-update"
)

// Release 版本
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset 静态资源
type Asset struct {
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// IsCompressedFile 返回是否是压缩文件的布尔值
func (a Asset) IsCompressedFile() bool {
	return a.ContentType == "application/zip" || a.ContentType == "application/x-gzip"
}

// ReleaseUpdater 版本更新器
type ReleaseUpdater struct {
	findAsset func(items []Asset) (idx int)
	client    *resty.Client
}

// NewReleaseUpdater 返回版本更新器实例
func NewReleaseUpdater(assetFinder func(items []Asset) (idx int)) *ReleaseUpdater {
	return &ReleaseUpdater{
		findAsset: assetFinder,
		client:    resty.New(),
	}
}

// CheckForUpdates 检查是否有更新
func (up ReleaseUpdater) CheckForUpdates(current *semver.Version, owner, repo string) (rel *Release, yes bool, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	var latest Release
	if _, err = up.client.R().SetHeader("Accept", "application/vnd.github.v3+json").SetResult(&latest).Get(url); err != nil {
		return nil, false, err
	}

	latestVersion, err := semver.NewVersion(latest.TagName)
	if err != nil {
		return nil, false, err
	}
	if latestVersion.GreaterThan(current) {
		return &latest, true, nil
	}
	return nil, false, nil
}

// Apply 更新指定版本
func (up ReleaseUpdater) Apply(rel *Release) (err error) {
	idx := up.findAsset(rel.Assets)
	if idx < 0 {
		return errs.ErrPackageNotFound
	}

	tmpDir, err := os.MkdirTemp("", strconv.FormatInt(time.Now().UnixNano(), 10))
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	url := rel.Assets[idx].BrowserDownloadURL
	srcFilename := filepath.Join(tmpDir, filepath.Base(url))
	dstFilename := srcFilename

	if _, err = up.client.R().SetOutput(srcFilename).Get(url); err != nil {
		return err
	}

	if rel.Assets[idx].IsCompressedFile() {
		if dstFilename, err = up.unarchive(srcFilename, tmpDir); err != nil {
			return err
		}
	}

	dstFile, err := os.Open(dstFilename)
	if err != nil {
		return nil
	}
	defer dstFile.Close()
	return update.Apply(dstFile, update.Options{})
}

// unarchive 解压缩至目标目录下并返回首个解压后的文件
func (up ReleaseUpdater) unarchive(srcFile, dstDir string) (dstFile string, err error) {
	if err = archiver.Unarchive(srcFile, dstDir); err != nil {
		return "", err
	}
	// 找到解压缩后的目标文件
	fis, _ := ioutil.ReadDir(dstDir)
	for _, fi := range fis {
		if strings.HasSuffix(srcFile, fi.Name()) {
			continue
		}
		return filepath.Join(dstDir, fi.Name()), nil
	}
	return "", nil
}
