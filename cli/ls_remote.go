package cli

import (
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/k0kubun/go-ansi"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/internal/collector"
	"github.com/voidint/g/internal/version"
)

const (
	stableChannel   = "stable"
	unstableChannel = "unstable"
	archivedChannel = "archived"
)

func listRemote(ctx *cli.Context) (err error) {
	channel := ctx.Args().First()
	if channel != "" && channel != stableChannel && channel != unstableChannel && channel != archivedChannel {
		return cli.ShowSubcommandHelp(ctx)
	}

	c, err := collector.NewCollector(strings.Split(os.Getenv(mirrorEnv), ",")...)
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	var vs []*version.Version
	switch channel {
	case stableChannel:
		vs, err = c.StableVersions()
	case unstableChannel:
		vs, err = c.UnstableVersions()
	case archivedChannel:
		vs, err = c.ArchivedVersions()
	default:
		vs, err = c.AllVersions()
	}
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	items := make([]*semver.Version, 0, len(vs))
	for _, ver := range vs {
		items = append(items, ver.SemanticVersion)
	}

	render(installed(), items, ansi.NewAnsiStdout())
	return nil
}
