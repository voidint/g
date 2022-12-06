package cli

import (
	"fmt"
	"os"

	"github.com/dixonwille/wmenu/v5"
	"github.com/urfave/cli/v2"
)

func selfUninstall(*cli.Context) (err error) {
	menu := wmenu.NewMenu("Are you sure you want to uninstall g?")
	menu.IsYesNo(wmenu.DefY)
	menu.Action(func(opts []wmenu.Opt) error {
		if opts[0].Value.(string) != "yes" {
			return nil
		}

		// Remove the g home directory and g binary files
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		rmPaths := []string{exePath}

		for {
			if binPath, err := os.Readlink(exePath); err == nil && binPath != exePath {
				rmPaths = append(rmPaths, binPath)
				exePath = binPath
			} else {
				break
			}
		}

		rmPaths = append(rmPaths, ghomeDir)

		var manRmPaths []string
		for i := range rmPaths {
			if err = os.RemoveAll(rmPaths[i]); err != nil {
				manRmPaths = append(manRmPaths, rmPaths[i])
			} else {
				fmt.Println("Remove", rmPaths[i])
			}
		}

		if len(manRmPaths) > 0 {
			fmt.Fprintln(os.Stderr, "Please manually remove the following files or directories:")
			for i := range manRmPaths {
				fmt.Fprintln(os.Stderr, manRmPaths[i])
			}
		}
		return nil
	})
	if err = menu.Run(); err != nil {
		return cli.Exit(errstring(err), 1)
	}
	return nil
}
