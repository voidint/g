package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
	"github.com/urfave/cli"
	"github.com/voidint/g/version"
)

func install(ctx *cli.Context) (err error) {
	vname := ctx.Args().First()
	if vname == "" {
		return cli.ShowSubcommandHelp(ctx)
	}
	targetV := filepath.Join(versionsDir, vname)

	// 检查版本是否已经安装
	if finfo, err := os.Stat(targetV); err == nil && finfo.IsDir() {
		return cli.NewExitError(fmt.Sprintf("[g] %q version has been installed.", vname), 1)
	}

	var url string
	if url = os.Getenv("G_MIRROR"); url == "" {
		url = version.DefaultURL
	}

	// 查找版本
	c, err := version.NewCollector(url)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	items, err := c.AllVersions()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	v, err := version.FindVersion(items, vname)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	// 查找版本下当前平台的安装包
	pkg, err := v.FindPackage(version.ArchiveKind, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	var ext string
	if runtime.GOOS == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}
	filename := filepath.Join(downloadsDir, fmt.Sprintf("go%s.%s-%s.%s", vname, runtime.GOOS, runtime.GOARCH, ext))

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// 本地不存在安装包，从远程下载并检查校验和。
		if _, err = pkg.Download(filename); err != nil {
			return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
		}
		if err = pkg.VerifyChecksum(filename); err != nil {
			return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
		}
	} else {
		// 本地存在安装包，检查校验和。
		if err = pkg.VerifyChecksum(filename); err != nil {
			return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
		}
	}
	// 解压安装包
	if err = archiver.Unarchive(filename, versionsDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	// 目录重命名
	if err = os.Rename(filepath.Join(versionsDir, "go"), targetV); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	// 重新建立软链接
	_ = os.Remove(goroot)

	if err := os.Symlink(targetV, goroot); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	fmt.Println("Installed successfully")
	return nil
}
