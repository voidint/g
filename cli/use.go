package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func use(ctx *cli.Context) (err error) {
	vname := ctx.Args().First()
	if vname == "" {
		return cli.ShowSubcommandHelp(ctx)
	}
	targetV := filepath.Join(versionsDir, vname)

	if finfo, err := os.Stat(targetV); err != nil || !finfo.IsDir() {
		return cli.NewExitError(fmt.Sprintf("[g] The %q version does not exist, please install it first.", vname), 1)
	}

	_ = os.Remove(goroot)

	if err := os.Symlink(targetV, goroot); err != nil {
		return cli.NewExitError(errstring(err), 1)
	}
	if output, err := exec.Command(filepath.Join(goroot, "bin", "go"), "version").Output(); err == nil {
		fmt.Print(string(output))
	}
	return nil
}
