package cli

import "github.com/urfave/cli"

var (
	commands = []cli.Command{

		{
			Name:      "list",
			ShortName: "ls",
			Usage:     "list versions installed",
			Subcommands: []cli.Command{
				{
					Name:   "known",
					Usage:  "list the versions of go available",
					Action: listKnown,
				},
			},
			Action: list,
		},
	}
)
