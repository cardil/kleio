package tasks

import (
	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

func Test(f *goyek.Flow) {
	f.Define(goyek.Task{
		Name:  "test",
		Usage: "Run tests",
		Deps: goyek.Deps{
			f.Define(Unit()),
		},
	})
}

func Unit() goyek.Task {
	return goyek.Task{
		Name:  "unit",
		Usage: "Runs unit tests for the project",
		Action: func(a *goyek.A) {
			cmd.Exec(a, "go test -v ./...")
		},
	}
}