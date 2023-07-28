package cli

import "github.com/urfave/cli/v2"

var (
	commands = []*cli.Command{
		{
			Name:      "ls",
			Aliases:   []string{"l"},
			Usage:     "List installed versions",
			UsageText: "g ls",
			Action:    list,
		},
		{
			Name:      "ls-remote",
			Aliases:   []string{"lr", "lsr"},
			Usage:     "List remote versions available for install",
			UsageText: "g ls-remote [stable|archived|unstable]",
			Action:    listRemote,
		},
		{
			Name:      "use",
			Usage:     "Switch to specified version",
			UsageText: "g use <version>",
			Action:    use,
		},
		{
			Name:      "install",
			Aliases:   []string{"i"},
			Usage:     "Download and install a version",
			UsageText: "g install <version>",
			Action:    install,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "nouse",
					Aliases: []string{"n"},
					Usage:   "Only install without using",
				},
			},
		},
		{
			Name:      "uninstall",
			Usage:     "Uninstall a version",
			UsageText: "g uninstall <version>",
			Action:    uninstall,
		},
		{
			Name:      "update",
			Usage:     "Download and install updates to g",
			UsageText: "g update",
			Action:    selfUpdate,
			Hidden:    true,
		},
		{
			Name:      "clean",
			Usage:     "Remove files from the package download directory",
			UsageText: "g clean",
			Action:    clean,
		},
		{
			Name:      "env",
			Usage:     "Show env variables of g",
			UsageText: "g env",
			Action:    showEnv,
		},
		{
			Name:  "self",
			Usage: "Modify g itself",
			Subcommands: []*cli.Command{
				{
					Name:      "update",
					Usage:     "Download and install updates to g",
					UsageText: "g self update",
					Action:    selfUpdate,
				},
				{
					Name:      "uninstall",
					Usage:     "Uninstall g",
					UsageText: "g self uninstall",
					Action:    selfUninstall,
				},
			},
		},
	}
)
