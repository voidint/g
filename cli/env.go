package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var envNames = []string{
	homeEnv,
	mirrorEnv,
	experimentalEnv,
}

func showEnv(ctx *cli.Context) (err error) {
	for _, e := range envNames {
		fmt.Printf("%s=%q\n", e, os.Getenv(e))
	}
	return nil
}
