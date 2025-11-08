package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/promo_handler"
)

func PromoRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	promoHandler := promo_handler.NewPromoHandler(app.Config, app.Usecases.PromoUsecase)

	promos := routerGroup.Group("/promos", middlewareMap.Auth, middlewareMap.TimeoutFast)
	{
		promos.GET("/", promoHandler.ListPromos)
		promos.POST("/", promoHandler.CreatePromo)
		promos.GET("/:id", promoHandler.PromoByID)
		promos.PUT("/:id", promoHandler.UpdatePromo)
		promos.DELETE("/:id", promoHandler.RemovePromo)
		promos.GET("/types", promoHandler.ListPromoTypes)
		promos.PUT("/status", promoHandler.SetStatusPromo)
	}
}
