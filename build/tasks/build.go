package tasks

import (
	"github.com/cardil/kleio/build/util/fs"
	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

func Build() goyek.Task {
	return goyek.Task{
		Name:  "build",
		Usage: "Builds the project",
		Action: func(a *goyek.A) {
			if err := fs.EnsureDir(fs.BuildOutputDir()); err != nil {
				a.Errorf("Failed to create output directory: %v", err)
				return
			}
			cmd.Exec(a, "go build -v -o build/output/collector ./cmd/collector")
		},
	}
}