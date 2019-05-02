package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli"
)

func use(ctx *cli.Context) (err error) {
	version := ctx.Args().First()
	if version == "" {
		return cli.ShowSubcommandHelp(ctx)
	}
	homeDir, _ := os.UserHomeDir()
	rootDir := filepath.Join(homeDir, ".g")
	versionDir := filepath.Join(rootDir, "versions", version)

	if finfo, err := os.Stat(versionDir); err != nil || !finfo.IsDir() {
		return cli.NewExitError(fmt.Sprintf("[g] The %q version does not exist, please install it first.", version), 1)
	}

	goDir := filepath.Join(rootDir, "go")
	_ = os.Remove(goDir)

	if err := os.Symlink(versionDir, goDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("[g] %s", err.Error()), 1)
	}
	if output, err := exec.Command(filepath.Join(goDir, "bin", "go"), "version").Output(); err == nil {
		fmt.Println(string(output))
	}
	return nil
}
