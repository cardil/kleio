package collector

import (
	"log/slog"

	"github.com/cardil/kleio/pkg/api"
	"github.com/cardil/kleio/pkg/collector"
	"github.com/cardil/kleio/pkg/server"
	"github.com/cardil/kleio/pkg/storage/inmem"
	"github.com/cardil/kleio/pkg/syslog"

	"github.com/wavesoftware/go-retcode"
)

type ExitFn func(retcode int)

func Serve() server.Server {
	st := inmem.NewStore()
	c := &collector.Collector{Storage: st}
	return server.Multi(
		syslog.Serve(c.Collect),
		api.Serve(st),
	)
}

func ServeOrDie(exit ExitFn) {
	err := Serve().Run()
	if err != nil {
		slog.Error("Bootstrap failure", "error", err)
		exit(retcode.Calc(err))
	}
}
