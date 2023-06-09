package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func showEnv(ctx *cli.Context) (err error) {
	var envVars = []string{experimentalEnv, homeEnv, mirrorEnv}

	for _, v := range envVars {
		fmt.Printf("%s = %s\n", v, os.Getenv(v))
	}
	return nil
}
