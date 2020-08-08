package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/k0kubun/go-ansi"
	"github.com/urfave/cli"
)

func list(ctx *cli.Context) (err error) {
	infos, err := ioutil.ReadDir(versionsDir)
	if err != nil || len(infos) <= 0 {
		fmt.Printf("No version installed yet\n\n")
		return nil
	}
	items := make([]*semver.Version, 0, len(infos))
	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		vname := infos[i].Name()
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
