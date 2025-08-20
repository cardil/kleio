package pipelines

import (
	"github.com/cardil/kleio/build/tasks"
	"github.com/goyek/goyek/v2"
)

func Default() *goyek.Flow {
	f := &goyek.Flow{}
	f.Define(tasks.Clean())
	tasks.Test(f)
	f.SetDefault(f.Define(tasks.Build()))
	return f
}