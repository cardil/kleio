package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		slog.Error("", "error", err)
	}

	if len(c.Errors) > 0 && !c.Writer.Written() {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	}
}
