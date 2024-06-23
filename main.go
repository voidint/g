package main

import (
	_ "github.com/voidint/g/collector/aliyun"
	_ "github.com/voidint/g/collector/ustc"

	"github.com/voidint/g/cli"
)

func main() {
	cli.Run()
}
