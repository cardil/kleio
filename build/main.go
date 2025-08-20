package main

import (
	"os"

	"github.com/cardil/kleio/build/boot"
	"github.com/cardil/kleio/build/pipelines"
	"github.com/cardil/kleio/build/util/fs"
	"github.com/goyek/goyek/v2"
)

func main() {
	if err := os.Chdir(fs.RootDir()); err != nil {
		panic(err)
	}
	goyek.DefaultFlow = pipelines.Default()
	boot.Main()
}