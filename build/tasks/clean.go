package tasks

import (
	"os"

	"github.com/cardil/kleio/build/util/fs"
	"github.com/goyek/goyek/v2"
)

func Clean() goyek.Task {
	return goyek.Task{
		Name:  "clean",
		Usage: "Clean build artifacts",
		Action: func(a *goyek.A) {
			outputDir := fs.BuildOutputDir()
			if err := os.RemoveAll(outputDir); err != nil {
				a.Errorf("Failed to clean output directory: %v", err)
				return
			}
			a.Logf("Cleaned %s", outputDir)
		},
	}
}