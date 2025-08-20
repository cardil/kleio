package pipelines

import (
	"github.com/cardil/kleio/build/tasks"
	"github.com/cardil/kleio/build/util/dotenv"
	"github.com/goyek/goyek/v2"
)

func Default() *goyek.Flow {
	f := &goyek.Flow{}
	f.UseExecutor(dotenv.Load)
	f.Define(tasks.Clean())
	f.Define(tasks.Format())
	f.Define(tasks.Vet())
	tasks.Test(f)
	f.SetDefault(f.Define(tasks.Build()))
	return f
}