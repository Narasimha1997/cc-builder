package main

import (
	"os"

	"github.com/Narasimha1997/cchelper/core"
)

func main() {
	args := os.Args
	rootPath := args[1]

	core.InitWatcher(rootPath)

}
