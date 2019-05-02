package cli

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// Run 运行g命令行
func Run() {
	app := cli.NewApp()
	app.Name = "g"
	app.Usage = "Golang version manager"
	app.Version = "0.1.0"
	app.Copyright = "Copyright (c) 2019, 2019, voidint. All rights reserved."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "voidint",
			Email: "voidint@126.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "D, debug",
			Usage: "enable debug mode",
			// Destination: &c.Debug,
		},
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
