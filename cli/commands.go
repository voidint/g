package cli

import "github.com/urfave/cli/v2"

var (
	commands = []*cli.Command{
		{
			Name:      "ls",
			Usage:     "List installed versions",
			UsageText: "g ls",
			Action:    list,
		},
		{
			Name:      "ls-remote",
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
			Usage:     "Download and install a version",
			UsageText: "g install <version>",
			Action:    install,
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
			Name:  "self",
			Usage: "Modify g itself",
			Subcommands: []*cli.Command{
				{
					Name:      "update",
					Usage:     "Download and install updates to g",
					UsageText: "g self update",
					Action:    selfUpdate,
				},
			},
		},
	}
)
