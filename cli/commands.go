package cli

import "github.com/urfave/cli"

var (
	commands = []cli.Command{
		{
			Name:   "ls",
			Usage:  "list installed versions",
			Action: list,
		},
		{
			Name:   "ls-remote",
			Usage:  "list remote versions available for install",
			Action: listRemote,
		},
	}
)
