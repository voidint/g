// Copyright (c) 2019 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var (
	commands = []*cli.Command{
		{
			Name:      "ls",
			Aliases:   []string{"l"},
			Usage:     "List installed versions",
			UsageText: "g ls",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "Output format. One of: [text|json]",
				},
			},
			Before: func(ctx *cli.Context) error {
				return validateLsFlag(ctx)
			},
			Action: list,
		},
		{
			Name:      "ls-remote",
			Aliases:   []string{"lr", "lsr"},
			Usage:     "List remote versions available for install",
			UsageText: "g ls-remote [stable|archived|unstable]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "Output format. One of: [text|json]",
				},
			},
			Before: func(ctx *cli.Context) error {
				return validateLsFlag(ctx)
			},
			Action: listRemote,
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
		{
			Name:      "mcp",
			Usage:     "Run in mcp server mode",
			UsageText: "g mcp",
			Action:    runMcpServer,
		},
	}
)

func validateLsFlag(ctx *cli.Context) error {
	if out := ctx.String("output"); out != "" && out != "json" && out != "text" {
		return cli.Exit(errstring(fmt.Errorf("unable to match a printer suitable for the output format %q, allowed formats are: [text|json]", out)), 1)
	}
	return nil
}
