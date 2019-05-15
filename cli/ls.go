package cli

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
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
		v, err := semver.NewVersion(infos[i].Name())
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
