package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/internal/build"
)

var (
	ghomeDir     string
	downloadsDir string
	versionsDir  string
	goroot       string
)

// Run 运行g命令行
func Run() {
	app := cli.NewApp()
	app.Name = "g"
	app.Usage = "Golang Version Manager"
	app.Version = build.Version()
	app.Copyright = fmt.Sprintf("Copyright (c) 2019-%d, voidint. All rights reserved.", time.Now().Year())
	app.Authors = []*cli.Author{
		{Name: "voidint", Email: "voidint@126.com"},
	}

	app.Before = func(ctx *cli.Context) (err error) {
		ghomeDir = ghome()
		goroot = filepath.Join(ghomeDir, "go")
		downloadsDir = filepath.Join(ghomeDir, "downloads")
		if err = os.MkdirAll(downloadsDir, 0755); err != nil {
			return err
		}
		versionsDir = filepath.Join(ghomeDir, "versions")
		return os.MkdirAll(versionsDir, 0755)
	}
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func init() {
	cli.AppHelpTemplate = fmt.Sprintf(`NAME:
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

 USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Commands}} command{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

 VERSION:
	%s{{end}}{{end}}{{if .Description}}

 DESCRIPTION:
	{{.Description}}{{end}}{{if len .Authors}}

 AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
	{{range $index, $author := .Authors}}{{if $index}}
	{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

 COMMANDS:{{range .VisibleCategories}}{{if .Name}}

	{{.Name}}:{{end}}{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

 GLOBAL OPTIONS:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

 COPYRIGHT:
	{{.Copyright}}{{end}}
`, build.ShortVersion)
}

const (
	experimentalEnv = "G_EXPERIMENTAL"
	homeEnv         = "G_HOME"
	mirrorEnv       = "G_MIRROR"
)

// ghome 返回g根目录
func ghome() (dir string) {
	if experimental := os.Getenv(experimentalEnv); experimental == "true" {
		if dir = os.Getenv(homeEnv); dir != "" {
			return dir
		}
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".g")
}

// inuse 返回当前的go版本号
func inuse(goroot string) (version string) {
	p, _ := os.Readlink(goroot)
	return filepath.Base(p)
}

// installed 返回当前的已经安装的go版本号
func installed() (versions map[string]bool) {
	dirs, err := os.ReadDir(versionsDir)
	if err != nil {
		return
	}

	inused := inuse(goroot)
	versions = make(map[string]bool, 0)
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		vname := d.Name()
		v, err := semversioned(vname)
		if err != nil || v == nil {
			continue
		}

		versions[vname] = (vname == inused)
	}

	return
}

// render 渲染go版本列表
func render(installed map[string]bool, items []*semver.Version, out io.Writer) {
	sort.Sort(semver.Collection(items))

	for i := range items {
		fields := strings.SplitN(items[i].String(), "-", 2)
		v := strings.TrimSuffix(strings.TrimSuffix(fields[0], ".0"), ".0")
		if len(fields) > 1 {
			v += fields[1]
		}
		if inused, found := installed[v]; found {
			if inused {
				color.New(color.FgGreen).Fprintf(out, "* %s\n", v)
			} else {
				color.New(color.FgGreen).Fprintf(out, "  %s\n", v)
			}
		} else {
			fmt.Fprintf(out, "  %s\n", v)
		}
	}
}

func semversioned(vname string) (*semver.Version, error) {
	var idx int
	if strings.Contains(vname, "alpha") {
		idx = strings.Index(vname, "alpha")

	} else if strings.Contains(vname, "beta") {
		idx = strings.Index(vname, "beta")

	} else if strings.Contains(vname, "rc") {
		idx = strings.Index(vname, "rc")
	}
	if idx > 0 {
		vname = vname[:idx] + "-" + vname[idx:]
	}
	return semver.NewVersion(vname)
}

// errstring 返回统一格式的错误信息
func errstring(err error) string {
	if err == nil {
		return ""
	}
	return wrapstring(err.Error())
}

func wrapstring(str string) string {
	if str == "" {
		return str
	}
	words := strings.Fields(str)
	if len(words) > 0 {
		words[0] = strings.Title(words[0])
	}
	return fmt.Sprintf("[g] %s", strings.Join(words, " "))
}
