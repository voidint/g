// Copyright (c) 2019 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cli

import (
	"fmt"
	"os"

	"github.com/dixonwille/wmenu/v5"
	"github.com/urfave/cli/v2"
)

func selfUninstall(*cli.Context) error {
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
			if binPath, e := os.Readlink(exePath); e == nil && binPath != exePath {
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
				fmt.Println("Removed", rmPaths[i])
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
	if err := menu.Run(); err != nil {
		return cli.Exit(errstring(err), 1)
	}
	return nil
}
