package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/banner_handler"

	"github.com/gin-gonic/gin"
)

func BannerRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {

	bannerHandler := banner_handler.NewBannerHandler(app.Usecases.BannerUsecase)

	banners := routerGroup.Group("/banners", middlewareMap.TimeoutFast)
	{
		banners.GET("/", middlewareMap.Auth, bannerHandler.ListBanners)
		banners.GET("/active", bannerHandler.ListActiveBanners)
		banners.GET("/:id", bannerHandler.DetailBanner)
		banners.POST("", middlewareMap.Auth, bannerHandler.CreateBanner, middlewareMap.TimeoutFile)
		banners.PUT("/:id", middlewareMap.Auth, bannerHandler.UpdateBanner, middlewareMap.TimeoutFile)
		banners.DELETE("/:id", middlewareMap.Auth, bannerHandler.RemoveBanner)
		banners.POST("/status", middlewareMap.Auth, bannerHandler.UpdateStatusBanner)
		banners.POST("/order", middlewareMap.Auth, bannerHandler.UpdateOrderBanner)
	}
}
