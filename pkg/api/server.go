package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/cardil/kleio/pkg/server"
	"github.com/cardil/kleio/pkg/storage"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

var ErrApiServer = errors.New("API Server failure")

const format = "20060102-150405"

func Serve(store storage.Storage) server.Server {
	a := &api{store: store}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		gin.Recovery(),
		sloggin.New(slog.Default()),
		ErrorHandler,
	)
	router.GET("/download", a.download)
	router.GET("/stats", a.stats)
	router.GET("/", a.home)

	port := 8080
	if sport, ok := os.LookupEnv("API_PORT"); ok {
		iport, err := strconv.Atoi(sport)
		if err == nil {
			port = iport
		}
	}
	bind := fmt.Sprint("0.0.0.0:", port)
	handler := router.Handler()
	srv := &http.Server{Addr: bind, Handler: handler}
	return &apiServ{server: srv, store: store}
}

type apiServ struct {
	server *http.Server
	store  storage.Storage
}

func (a *apiServ) Run() (err error) {
	slog.Info("Starting API server", "bind", a.server.Addr)
	err = a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrApiServer, err)
	}
	return
}

func (a *apiServ) Close() (err error) {
	ctx := context.Background()
	err = a.server.Shutdown(ctx)
	err = errors.Join(err, a.store.Close())
	return
}

type api struct {
	store storage.Storage
}

func (a *api) home(c *gin.Context) {
	c.Redirect(http.StatusFound, "/stats")
}
