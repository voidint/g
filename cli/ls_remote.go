package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/urfave/cli"
	"github.com/voidint/g/version"
)

const (
	stableChannel   = "stable"
	archivedChannel = "archived"
)

func listRemote(ctx *cli.Context) (err error) {
	channel := ctx.Args().First()
	if channel != "" && channel != stableChannel && channel != archivedChannel {
		return cli.ShowSubcommandHelp(ctx)
	}

	var url string
	if url = os.Getenv("G_MIRROR"); url == "" {
		url = version.DefaultURL
	}

	c, err := version.NewCollector(url)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}

	var vs []*version.Version
	switch channel {
	case stableChannel:
		vs, err = c.StableVersions()
	case archivedChannel:
		vs, err = c.ArchivedVersions()
	default:
		vs, err = c.AllVersions()
	}
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}

	items := make([]*semver.Version, 0, len(vs))
	for i := range vs {
		v, err := semver.NewVersion(vs[i].Name)
		if err != nil || v == nil {
			continue
		}
		items = append(items, v)
	}
	sort.Sort(semver.Collection(items))

	for i := range items {
		fmt.Println(strings.TrimSuffix(strings.TrimSuffix(items[i].String(), ".0"), ".0"))
	}
	return nil
}
