package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *api) download(context *gin.Context) {
	slog.Info("Download")
	reader, err := a.store.Download().ZipReader()
	if err != nil {
		_ = context.Error(err)
		return
	}
	contentType := "application/zip"
	now := time.Now()
	filename := fmt.Sprintf("logs-%s.zip", now.Format(format))
	headers := map[string]string{
		"Content-Disposition": `attachment; filename="` + filename + `"`,
	}
	context.DataFromReader(http.StatusOK, reader.Size, contentType, reader, headers)
	if err = reader.Close(); err != nil {
		_ = context.Error(err)
	}
}
