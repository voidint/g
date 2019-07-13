package cli

import (
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/urfave/cli"
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
	if url = os.Getenv("G_MIRROR"); url == "" {
		url = version.DefaultURL
	}

	c, err := version.NewCollector(url)
	if err != nil {
		return cli.NewExitError(errstring(err), 1)
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
		return cli.NewExitError(errstring(err), 1)
	}

	items := make([]*semver.Version, 0, len(vs))
	for i := range vs {
		vname := vs[i].Name
		idx := strings.Index(vname, "beta")
		if idx > 0 {
			vname = vname[:idx] + "-" + vname[idx:]
		}
		v, err := semver.NewVersion(vname)
		if err != nil || v == nil {
			continue
		}
		items = append(items, v)
	}

	render(inuse(goroot), items, os.Stdout)
	return nil
}
