package cli

import (
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/k0kubun/go-ansi"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/collector"
	"github.com/voidint/g/version"
)

const (
	stableChannel   = "stable"
	unstableChannel = "unstable"
	archivedChannel = "archived"
)

func listRemote(ctx *cli.Context) (err error) {
	vname := ctx.Args().First()

	var cs *semver.Constraints
	if vname != "" && vname != stableChannel && vname != unstableChannel && vname != archivedChannel && vname != version.Latest {
		if cs, err = semver.NewConstraint(vname); err != nil {
			return cli.Exit(errstring(err), 1)
		}
	}

	c, err := collector.NewCollector(strings.Split(os.Getenv(mirrorEnv), mirrorSep)...)
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	var vs []*version.Version
	switch vname {
	case stableChannel:
		vs, err = c.StableVersions()
	case unstableChannel:
		vs, err = c.UnstableVersions()
	case archivedChannel:
		vs, err = c.ArchivedVersions()
	default:
		vs, err = c.AllVersions()
		if err == nil {
			if vname == version.Latest {
				vs = []*version.Version{vs[len(vs)-1]}
			}

			if vname != "" && vname != version.Latest {
				var newVs []*version.Version
				for _, v := range vs {
					if v.MatchConstraint(cs) {
						newVs = append(newVs, v)
					}
				}
				vs = newVs
			}
		}

	}
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	var renderMode uint8
	switch ctx.String("output") {
	case "json":
		renderMode = jsonMode
	default:
		renderMode = textMode
	}

	render(renderMode, installed(), vs, ansi.NewAnsiStdout())
	return nil
}
