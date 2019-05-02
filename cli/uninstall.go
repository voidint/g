package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func uninstall(ctx *cli.Context) (err error) {
	version := ctx.Args().First()
	if version == "" {
		return cli.ShowSubcommandHelp(ctx)
	}
	homeDir, _ := os.UserHomeDir()
	rootDir := filepath.Join(homeDir, ".g")
	versionDir := filepath.Join(rootDir, "versions", version)

	if finfo, err := os.Stat(versionDir); err != nil || !finfo.IsDir() {
		return cli.NewExitError(fmt.Sprintf("[g] %q version is not installed.", version), 1)
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] Uninstall failed ==> %s", err.Error()), 1)
	}
	fmt.Println("Uninstall successfully")
	return nil
}
