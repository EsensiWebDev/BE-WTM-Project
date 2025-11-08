package middleware

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"time"
)

func TimeoutMiddleware(d time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(d),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(408, gin.H{"error": "Request Timeout"})
		}),
	)
}
