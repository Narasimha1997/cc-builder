package main

import (
	"os"

	"github.com/Narasimha1997/ccbuilder/core"
)

func main() {
	args := os.Args
	rootPath := args[1]

	core.InitWatcher(rootPath)

}
