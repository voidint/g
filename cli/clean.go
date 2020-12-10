package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func clean(ctx *cli.Context) (err error) {
	entries, err := ioutil.ReadDir(downloadsDir)
	if err != nil {
		return cli.NewExitError(errstring(err), 1)
	}

	for i := range entries {
		if err = os.RemoveAll(filepath.Join(downloadsDir, entries[i].Name())); err == nil {
			fmt.Println("Remove", entries[i].Name())
		}
	}
	return nil
}
