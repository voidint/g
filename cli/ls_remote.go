package cli

import (
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/k0kubun/go-ansi"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/collector"
	"github.com/voidint/g/collector/official"
	"github.com/voidint/g/version"
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

	var url string
	if url = os.Getenv(mirrorEnv); url == "" {
		url = official.DefaultDownloadPageURL
	}

	c, err := collector.NewCollector(url)
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
	for i := range vs {
		vname := vs[i].Name
		var idx int
		if strings.Contains(vname, "alpha") {
			idx = strings.Index(vname, "alpha")

		} else if strings.Contains(vname, "beta") {
			idx = strings.Index(vname, "beta")

		} else if strings.Contains(vname, "rc") {
			idx = strings.Index(vname, "rc")
		}
		if idx > 0 {
			vname = vname[:idx] + "-" + vname[idx:]
		}
		v, err := semver.NewVersion(vname)
		if err != nil || v == nil {
			continue
		}
		items = append(items, v)
	}

	render(inuse(goroot), items, ansi.NewAnsiStdout())
	return nil
}
