package cli

import "github.com/urfave/cli"

var (
	commands = []cli.Command{
		{
			Name:   "ls",
			Usage:  "List installed versions",
			Action: list,
		},
		{
			Name:      "ls-remote",
			Usage:     "List remote versions available for install",
			UsageText: "g ls-remote [stable|archived]",
			Action:    listRemote,
		},
		{
			Name:      "use",
			Usage:     "Switch to specified version",
			UsageText: "g use <version>",
			Action:    use,
		},
	}
)
