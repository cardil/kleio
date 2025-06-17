package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) stats(context *gin.Context) {
	slog.Info("Stats")
	stats := a.store.Stats()
	context.JSON(http.StatusOK, stats)
}
