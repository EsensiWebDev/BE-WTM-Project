package main

import (
	"context"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"wtm-backend/docs"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/infrastructure/router"
	"wtm-backend/pkg/logger"
)

// @title WTM Backend Service
// @description This is a group of API for WTM Service
// @version 1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	// Inisialisasi logger
	logger.InitLogger()
	logger.Info(ctx, "Logger initialized")

	// Inisialisasi aplikasi
	app := bootstrap.NewApplication()

	// Jalankan server
	port := app.Config.ServerPort
	logger.Info(ctx, "Server started on port "+port)

	// Setup router Gin
	r := router.SetupRouter(app)

	// Swagger endpoint
	if !app.Config.IsProduction() {
		go func() {
			log.Println("Starting pprof on :6060")
			if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
				log.Println("pprof error:", err)
			}
		}()
		host := app.Config.Host
		schema := []string{"https"}
		if host == "" {
			host = "localhost:" + port // fallback default
		}
		logger.Info(ctx, "Server started on host "+host)

		if strings.Contains(host, ":") {
			schema = []string{"http"}
		}

		docs.SwaggerInfo.Host = host
		docs.SwaggerInfo.BasePath = "/api"
		docs.SwaggerInfo.Schemes = schema
		docs.SwaggerInfo.Version = ""
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		logger.Info(ctx, "Swagger disabled in production")
	}

	err := r.Run("0.0.0.0:" + port)
	if err != nil {
		logger.Fatal("Failed to start server", err.Error())
	}
}
