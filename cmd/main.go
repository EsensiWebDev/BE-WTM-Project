package main

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

	// Initialize core components
	if err := initializeApplication(ctx); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Setup application
	app := bootstrap.NewApplication()
	setupRouter := router.SetupRouter(app)

	// Configure Swagger for non-production
	configureSwagger(app, setupRouter)

	// Configure pprof for debugging
	configurePprof(app)

	// Start server with graceful shutdown
	startServer(ctx, app, setupRouter)
}

func initializeApplication(ctx context.Context) error {
	logger.InitLogger()
	logger.Info(ctx, "Logger initialized successfully")
	return nil
}

func configureSwagger(app *bootstrap.Application, router *gin.Engine) {
	if app.Config.IsProduction() {
		logger.Info(context.Background(), "Swagger disabled in production environment")
		return
	}

	// Configure Swagger
	host := getSwaggerHost(app)
	schemes := getSwaggerSchemes(host)

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = schemes
	docs.SwaggerInfo.Version = "1.0" // ⭐ Ditambahkan version

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info(context.Background(), "Swagger documentation available at /swagger/index.html")
}

func configurePprof(app *bootstrap.Application) {
	if !app.Config.IsProduction() {
		go func() {
			logger.Info(context.Background(), "Starting pprof debug server on :6060")
			if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
				logger.Error(context.Background(), "Pprof server error: "+err.Error())
			}
		}()
	}
}

func getSwaggerHost(app *bootstrap.Application) string {
	if app.Config.Host != "" {
		host := strings.TrimPrefix(app.Config.Host, "http://")
		host = strings.TrimPrefix(host, "https://")
		return host
	}
	return "localhost:" + app.Config.ServerPort
}

func getSwaggerSchemes(host string) []string {
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
		return []string{"http"}
	}
	// kalau ada port → http
	if strings.Contains(host, ":") {
		return []string{"http"}
	}
	return []string{"https"}
}

func startServer(ctx context.Context, app *bootstrap.Application, router *gin.Engine) {
	port := app.Config.ServerPort
	address := "0.0.0.0:" + port

	logger.Info(ctx, "Starting server on "+address)

	// Create server with timeouts for production readiness
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine to handle graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "Server failed to start: "+err.Error())
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "Server forced to shutdown: "+err.Error())
	}

	logger.Info(ctx, "Server exited properly")
}
