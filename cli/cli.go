package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/voidint/g/build"
)

// Run 运行g命令行
func Run() {
	app := cli.NewApp()
	app.Name = "g"
	app.Usage = "Golang version manager"
	app.Version = build.Version()
	app.Copyright = "Copyright (c) 2019, 2019, voidint. All rights reserved."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "voidint",
			Email: "voidint@126.com",
		},
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[g] %s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	cli.AppHelpTemplate = fmt.Sprintf(`NAME:
	{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
 
 USAGE:
	{{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
 
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
