package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/banner_handler"
)

func BannerRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {

	bannerHandler := banner_handler.NewBannerHandler(app.Usecases.BannerUsecase)

	banners := routerGroup.Group("/banners", middlewareMap.Auth, middlewareMap.TimeoutFast)
	{
		banners.GET("/", bannerHandler.ListBanners)
		banners.POST("", bannerHandler.CreateBanner, middlewareMap.TimeoutFile)
		banners.GET("/:id", bannerHandler.DetailBanner)
		banners.PUT("/:id", bannerHandler.UpdateBanner, middlewareMap.TimeoutFile)
		banners.DELETE("/:id", bannerHandler.RemoveBanner)
		banners.POST("/status", bannerHandler.UpdateStatusBanner)
		banners.POST("/order", bannerHandler.UpdateOrderBanner)
	}
}
