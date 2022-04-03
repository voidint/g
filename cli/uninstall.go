package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func uninstall(ctx *cli.Context) (err error) {
	vname := ctx.Args().First()
	if vname == "" {
		return cli.ShowSubcommandHelp(ctx)
	}
	targetV := filepath.Join(versionsDir, vname)

	if finfo, err := os.Stat(targetV); err != nil || !finfo.IsDir() {
		return cli.Exit(fmt.Sprintf("[g] %q version is not installed", vname), 1)
	}

	if err = os.RemoveAll(targetV); err != nil {
		return cli.Exit(wrapstring(fmt.Sprintf("Uninstall failed: %s", err.Error())), 1)
	}
	fmt.Printf("Uninstalled go%s\n", vname)
	return nil
}
