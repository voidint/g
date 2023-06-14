package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func clean(*cli.Context) (err error) {
	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	for i := range entries {
		if err = os.RemoveAll(filepath.Join(downloadsDir, entries[i].Name())); err == nil {
			fmt.Println("Remove", entries[i].Name())
		}
	}
	return nil
}
