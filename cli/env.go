package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var envNames = []string{
	homeEnv,
	mirrorEnv,
	experimentalEnv,
}

func showEnv(ctx *cli.Context) (err error) {
	printEnv(os.Stdout, envNames)
	return nil
}

func printEnv(w io.Writer, envNames []string) {
	for _, eName := range envNames {
		_, _ = fmt.Fprintf(w, "%s=%q\n", eName, os.Getenv(eName))
	}
}
