package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func showEnv(ctx *cli.Context) (err error) {
	var envVars = []string{experimentalEnv, homeEnv, mirrorEnv}

	for _, e := range envVars {
		v, set := os.LookupEnv(e)
		if !set {
			color.New(color.FgYellow).Fprintf(os.Stdout, "%s=(UNSET)\n", e)
		} else {
			fmt.Printf("%s=%s\n", e, v)
		}
	}
	return nil
}
