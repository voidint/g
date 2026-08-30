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
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/voidint/g/version"
)

func prune(ctx *cli.Context) error {
	items, err := listLocalVersions(versionsDir)
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}
	if len(items) == 0 {
		fmt.Println("No version installed yet")
		return nil
	}

	doomed := superseded(items, inuse(goroot))
	if len(doomed) == 0 {
		fmt.Println("Nothing to uninstall")
		return nil
	}

	dryRun := ctx.Bool("dry-run")
	for _, v := range doomed {
		if dryRun {
			fmt.Printf("Would uninstall %s\n", v.Name())
			continue
		}
		if err := os.RemoveAll(filepath.Join(versionsDir, v.Name())); err != nil {
			return cli.Exit(wrapstring(fmt.Sprintf("Uninstall %q failed: %s", v.Name(), err.Error())), 1)
		}
		fmt.Printf("Uninstalled %s\n", v.Name())
	}
	return nil
}

// superseded returns the versions that have a newer version in the same minor
// series, excluding the version currently in use. It assumes items are sorted
// in ascending order (see listLocalVersions).
func superseded(items []*version.Version, inuse string) (doomed []*version.Version) {
	for i, v := range items {
		if v.Name() == inuse {
			continue // never remove the version in use
		}
		// Same-series versions are adjacent in ascending order, so a version
		// is superseded when its successor belongs to the same minor series.
		if i+1 < len(items) && sameMinorSeries(v, items[i+1]) {
			doomed = append(doomed, v)
		}
	}
	return doomed
}

// sameMinorSeries reports whether two versions belong to the same major.minor series.
func sameMinorSeries(a, b *version.Version) bool {
	return a.Major() == b.Major() && a.Minor() == b.Minor()
}
