package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func list(c *cli.Context) (err error) {
	homeDir, _ := os.UserHomeDir()

	infos, err := ioutil.ReadDir(filepath.Join(homeDir, ".g", "versions"))
	if err != nil {
		fmt.Printf("No version installed yet\n\n")
		return nil
	}
	// TODO 排序

	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		fmt.Println(infos[i].Name())
	}
	fmt.Println()
	return nil
}
