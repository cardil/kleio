package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/cardil/qe-clusterlogging/pkg/server"
	"github.com/cardil/qe-clusterlogging/pkg/storage"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

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
	return &apiServ{server: srv}
}

type apiServ struct {
	server *http.Server
}

func (h *apiServ) Run() error {
	return h.server.ListenAndServe()
}

func (h *apiServ) Kill() error {
	return h.server.Shutdown(context.Background())
}

type api struct {
	store storage.Storage
}

func (a *api) home(c *gin.Context) {
	c.Redirect(http.StatusFound, "/stats")
}
