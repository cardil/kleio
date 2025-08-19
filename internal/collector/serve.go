package collector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cardil/kleio/pkg/api"
	"github.com/cardil/kleio/pkg/collector"
	"github.com/cardil/kleio/pkg/server"
	"github.com/cardil/kleio/pkg/storage/facade"
	"github.com/cardil/kleio/pkg/syslog"

	"github.com/wavesoftware/go-retcode"
)

var ErrBootstrap = errors.New("bootstrap failure")

type ExitFn func(retcode int)

func NewServer() (server.Server, error) {
	ctx := context.Background()
	st, err := facade.NewStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrBootstrap, err)
	}
	c := &collector.Collector{Storage: st}
	return server.Multi(
		syslog.Serve(c.Collect),
		api.Serve(st),
	), nil
}

func Serve() (err error) {
	var srv server.Server
	if srv, err = NewServer(); err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, srv.Close())
	}()
	err = errors.Join(err, srv.Run())
	return
}

func ServeOrDie(exit ExitFn) {
	err := Serve()
	if err != nil {
		slog.Error("Failure", "error", err)
		exit(retcode.Calc(err))
	}
}
