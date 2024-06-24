package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/build"
	"github.com/voidint/g/version"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		if err = os.MkdirAll(downloadsDir, 0750); err != nil {
			return err
		}
		versionsDir = filepath.Join(ghomeDir, "versions")
		return os.MkdirAll(versionsDir, 0750)
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

const (
	// mirror site separator
	mirrorSep = ","
)

// ghome 返回g根目录
func ghome() (dir string) {
	if experimental := os.Getenv(experimentalEnv); strings.EqualFold(experimental, "true") {
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
		versions[vname] = (vname == inused)
	}

	return
}

type versionOut struct {
	Version   string            `json:"version"`
	InUse     bool              `json:"inUse"`
	Installed bool              `json:"installed"`
	Packages  []version.Package `json:"packages,omitempty"`
}

const (
	textMode = 0
	jsonMode = 1
)

// render 渲染go版本列表
func render(mode uint8, installed map[string]bool, items []*version.Version, out io.Writer) {
	switch mode {
	case jsonMode:
		vs := make([]versionOut, 0, len(items))

		for _, item := range items {
			vo := versionOut{
				Version:  item.Name(),
				Packages: item.Packages(),
			}
			if inuse, found := installed[item.Name()]; found {
				vo.InUse = inuse
				vo.Installed = found
			}
			vs = append(vs, vo)
		}

		enc := json.NewEncoder(out)
		enc.SetIndent("", "    ")
		_ = enc.Encode(&vs)

	default:
		for _, item := range items {
			if inused, found := installed[item.Name()]; found {
				if inused {
					_, _ = color.New(color.FgGreen).Fprintf(out, "* %s\n", item.Name())
				} else {
					_, _ = color.New(color.FgGreen).Fprintf(out, "  %s\n", item.Name())
				}
			} else {
				_, _ = fmt.Fprintf(out, "  %s\n", item.Name())
			}
		}
	}
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
		words[0] = cases.Title(language.English).String(words[0])
	}
	return fmt.Sprintf("[g] %s", strings.Join(words, " "))
}
