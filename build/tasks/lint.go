package tasks

import (
	"github.com/goyek/goyek/v2"
	"github.com/goyek/x/cmd"
)

func Format() goyek.Task {
	return goyek.Task{
		Name:  "format",
		Usage: "Format Go code",
		Action: func(a *goyek.A) {
			cmd.Exec(a, "go fmt ./...")
		},
	}
}

func Vet() goyek.Task {
	return goyek.Task{
		Name:  "vet",
		Usage: "Run go vet on the code",
		Action: func(a *goyek.A) {
			cmd.Exec(a, "go vet ./...")
		},
	}
}