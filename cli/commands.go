package cli

import "github.com/urfave/cli"

var (
	commands = []cli.Command{
		{
			Name:      "ls",
			Usage:     "List installed versions",
			UsageText: "g ls",
			Action:    list,
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
	}
)
