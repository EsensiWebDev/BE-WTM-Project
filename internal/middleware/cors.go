package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: m.allowedOrigins,

		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},

		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // ⬅️ penting untuk cookie / Authorization
		MaxAge:           m.maxAgeCors,
	}

	return cors.New(config)
}
