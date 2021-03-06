package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health is simple keep-alive/ping handler
func Health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	}
}
