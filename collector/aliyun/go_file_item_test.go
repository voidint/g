package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidint/g/version"
)

func Test_getGoVersion(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected string
	}{
		{
			In: &goFileItem{
				FileName: "go1.18beta1.darwin-amd64.pkg",
				URL:      "https://mirrors.aliyun.com/golang/go1.18beta1.darwin-amd64.pkg",
				Size:     "136.9 MB",
			},
			Expected: "1.18beta1",
		},
		{
			In: &goFileItem{
				FileName: "go1.18beta2.freebsd-386.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.18beta2.freebsd-386.tar.gz",
				Size:     "107.2 MB",
			},
			Expected: "1.18beta2",
		},
		{
			In: &goFileItem{
				FileName: "go1.18rc1.darwin-amd64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.18rc1.darwin-amd64.tar.gz",
				Size:     "107.2 MB",
			},
			Expected: "1.18rc1",
		},
		{
			In: &goFileItem{
				FileName: "go1.18.windows-arm64.zip",
				URL:      "https://mirrors.aliyun.com/golang/go1.18.windows-arm64.zip",
				Size:     "118.0 MB",
			},
			Expected: "1.18",
		},
		{
			In: &goFileItem{
				FileName: "go1.18.1.linux-386.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.18.1.linux-386.tar.gz",
				Size:     "107.6 MB",
			},
			Expected: "1.18.1",
		},
		{
			In: &goFileItem{
				FileName: "go1.18.1.src.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.18.1.src.tar.gz",
				Size:     "21.8 MB",
			},
			Expected: "1.18.1",
		},
	}
	t.Run("从文件名中获取go版本号", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.getGoVersion())
		}
	})
}

func Test_isSHA256File(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected bool
	}{
		{
			In: &goFileItem{
				FileName: "go1.18beta1.darwin-amd64.pkg",
				URL:      "https://mirrors.aliyun.com/golang/go1.18beta1.darwin-amd64.pkg",
				Size:     "136.9 MB",
			},
			Expected: false,
		},
		{
			In: &goFileItem{
				FileName: "go1.4-bootstrap-20170518.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.4-bootstrap-20170518.tar.gz.sha256",
				Size:     "64.0 B",
			},
			Expected: true,
		},
		{
			In: &goFileItem{
				FileName: "go1.4-bootstrap-20170518.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.4-bootstrap-20170518.tar.gz.sha256",
				Size:     "64.0 B",
			},
			Expected: true,
		},
		{
			In: &goFileItem{
				FileName: "go1.5.4.darwin-amd64.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.5.4.darwin-amd64.tar.gz.sha256",
				Size:     "64.0 B",
			},
			Expected: true,
		},
	}
	t.Run("判断是否是SHA256校验和文件", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.isSHA256File())
		}
	})
}

func Test_isPackageFile(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected bool
	}{
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg",
			},
			Expected: true,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg.sha256",
			},
			Expected: false,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-amd64.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz.sha256",
			},
			Expected: false,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz",
			},
			Expected: true,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz.asc",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz.asc",
			},
			Expected: false,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.msi",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.msi",
			},
			Expected: true,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.zip",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.zip",
			},
			Expected: true,
		},
	}
	t.Run("判断是否是go安装包文件", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.isPackageFile())
		}
	})
}

func Test_getKind(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected version.PackageKind
	}{
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg",
			},
			Expected: version.InstallerKind,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg.sha256",
			},
			Expected: "Unknown",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-amd64.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz.sha256",
			},
			Expected: "Unknown",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz",
			},
			Expected: version.ArchiveKind,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz.asc",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz.asc",
			},
			Expected: "Unknown",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.msi",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.msi",
			},
			Expected: version.InstallerKind,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.zip",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.zip",
			},
			Expected: version.ArchiveKind,
		},
		{
			In: &goFileItem{
				FileName: "go1.10.2.src.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.2.src.tar.gz",
			},
			Expected: version.SourceKind,
		},
	}
	t.Run("从文件名中获取文件类型", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.getKind())
		}
	})
}

func Test_getOS(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected string
	}{
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg",
			},
			Expected: "macOS",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.darwin-amd64.pkg.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.darwin-amd64.pkg.sha256",
			},
			Expected: "macOS",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-amd64.tar.gz.sha256",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz.sha256",
			},
			Expected: "Linux",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz",
			},
			Expected: "Linux",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.linux-arm64.tar.gz.asc",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.linux-arm64.tar.gz.asc",
			},
			Expected: "Linux",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.msi",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.msi",
			},
			Expected: "Windows",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.windows-386.zip",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.windows-386.zip",
			},
			Expected: "Windows",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.2.src.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.2.src.tar.gz",
			},
			Expected: "",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.3.freebsd-386.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.3.freebsd-386.tar.gz",
			},
			Expected: "FreeBSD",
		},
	}
	t.Run("从文件名中获取操作系统", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.getOS())
		}
	})
}

func Test_getArch(t *testing.T) {
	items := []*struct {
		In       *goFileItem
		Expected string
	}{
		{
			In: &goFileItem{
				FileName: "go1.10.2.src.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.2.src.tar.gz",
			},
			Expected: "",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.freebsd-386.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.freebsd-386.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "x86",
		},
		{
			In: &goFileItem{
				FileName: "go1.10.1.freebsd-amd64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.10.1.freebsd-amd64.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "x86-64",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.0.freebsd-arm.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.0.freebsd-arm.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "ARMv6",
		},
		{
			In: &goFileItem{
				FileName: "go1.20rc3.windows-arm64.msi",
				URL:      "https://mirrors.aliyun.com/golang/go1.20rc3.windows-arm64.msi?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "ARM64",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.0.linux-armv6l.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.0.linux-armv6l.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "ARMv6",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.0.aix-ppc64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.0.aix-ppc64.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "ppc64",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.0.linux-ppc64le.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.0.linux-ppc64le.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "ppc64le",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.1.linux-mips.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.1.linux-mips.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "mips",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.2.linux-mipsle.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.2.linux-mipsle.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "mipsle",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.3.linux-mips64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.3.linux-mips64.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "mips64",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.3.linux-mips64le.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.3.linux-mips64le.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "mips64le",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.3.linux-s390x.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.3.linux-s390x.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "s390x",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.4.freebsd-riscv64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.4.freebsd-riscv64.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "riscv64",
		},
		{
			In: &goFileItem{
				FileName: "go1.21.4.linux-loong64.tar.gz",
				URL:      "https://mirrors.aliyun.com/golang/go1.21.4.linux-loong64.tar.gz?spm=a2c6h.25603864.0.0.a6b07c45KMj71H",
			},
			Expected: "loong64",
		},
	}
	t.Run("从文件名中获取架构", func(t *testing.T) {
		for _, item := range items {
			assert.Equal(t, item.Expected, item.In.getArch())
		}
	})
}
