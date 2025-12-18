package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/middleware"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MiddlewareMap struct {
	Auth              gin.HandlerFunc
	TimeoutFast       gin.HandlerFunc
	TimeoutSlow       gin.HandlerFunc
	TimeoutFile       gin.HandlerFunc
	RequirePermission func(required string) gin.HandlerFunc
	RequireRole       func(required string) gin.HandlerFunc
}

func SetupRouter(app *bootstrap.Application) *gin.Engine {
	route := gin.Default()

	// Endpoint metrik untuk Prometheus
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//route.Use(app.Middleware.CORSMiddleware())
	route.Use(gin.Logger())
	route.Use(gin.Recovery())
	route.Use(middleware.TraceIDMiddleware())

	middlewareMap := MiddlewareMap{
		Auth:              app.Middleware.AuthMiddleware(),
		TimeoutFast:       middleware.TimeoutMiddleware(app.Config.DurationCtxTOFast),
		TimeoutSlow:       middleware.TimeoutMiddleware(app.Config.DurationCtxTOSlow),
		TimeoutFile:       middleware.TimeoutMiddleware(app.Config.DurationCtxTOFile),
		RequirePermission: app.Middleware.RequirePermission,
		RequireRole:       app.Middleware.RequireRole,
	}

	api := route.Group("api")
	{
		api.GET("/ping", PingHandler)
		AuthRoutes(app, middlewareMap, api)
		UserRoutes(app, middlewareMap, api)
		PromoRoutes(app, middlewareMap, api)
		HotelRoute(app, middlewareMap, api)
		BannerRoutes(app, middlewareMap, api)
		PromoGroupRoutes(app, middlewareMap, api)
		BookingRoute(app, middlewareMap, api)
		ReportRouter(app, middlewareMap, api)
		EmailRoute(app, middlewareMap, api)
		NotificationRouter(app, middlewareMap, api)
		CurrencyRoutes(app, middlewareMap, api)
		FileRouter(app, api)
	}
	return route
}

// PingHandler memeriksa apakah server berjalan.
// @Summary Ping API
// @Description Cek apakah server berjalan
// @Tags Health Check
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Ping endpoint accessed")
	c.JSON(200, gin.H{"message": "pong"})
}
