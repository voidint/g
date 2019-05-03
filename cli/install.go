package cli

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
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
	homeDir, _ := os.UserHomeDir()
	rootDir := filepath.Join(homeDir, ".g")
	versionDir := filepath.Join(rootDir, "versions", vname)

	// 检查版本是否已经安装
	if finfo, err := os.Stat(versionDir); err == nil && finfo.IsDir() {
		return cli.NewExitError(fmt.Sprintf("[g] %q version has been installed.", vname), 1)
	}

	// 查找版本
	c, err := version.NewCollector("https://golang.google.cn/dl/")
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

	filename := filepath.Join(homeDir, ".g", "downloads", fmt.Sprintf("go%s.%s-%s.tar.gz", vname, runtime.GOOS, runtime.GOARCH))

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// 本地不存在安装包，从远程下载并检查校验和。
		if _, err = v.Download(version.ArchiveKind, runtime.GOOS, runtime.GOARCH, filename); err != nil {
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
	if err = archiver.Unarchive(filename, filepath.Join(rootDir, "versions")); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	// 目录重命名
	if err = os.Rename(filepath.Join(rootDir, "versions", "go"), versionDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	// 重新建立软链接
	goDir := filepath.Join(rootDir, "go")
	_ = os.Remove(goDir)

	if err := os.Symlink(versionDir, goDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	fmt.Println("installed successfully")
	return nil
}

// unarchive 解压tar.gz文件
func unarchive(targzFile, dstDir string) (err error) {
	if err = os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(targzFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		filename := filepath.Join(dstDir, hdr.Name)
		os.MkdirAll(filepath.Dir(filename), 0755)
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		io.Copy(f, tr)
		_ = f.Close()
	}
	return nil
}
