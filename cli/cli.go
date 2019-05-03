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
