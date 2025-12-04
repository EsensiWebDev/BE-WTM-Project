package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/promo_handler"

	"github.com/gin-gonic/gin"
)

func PromoRoutes(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	promoHandler := promo_handler.NewPromoHandler(app.Config, app.Usecases.PromoUsecase)

	promos := routerGroup.Group("/promos", mm.Auth, mm.TimeoutFast)
	{
		promos.GET("/", mm.RequirePermission("promo:view"), promoHandler.ListPromos)
		promos.GET("/agent", promoHandler.ListPromoForAgent)
		promos.POST("/", mm.RequirePermission("promo:create"), promoHandler.CreatePromo)
		promos.GET("/:id", mm.RequirePermission("promo:view"), promoHandler.PromoByID)
		promos.PUT("/:id", mm.RequirePermission("promo:edit"), promoHandler.UpdatePromo)
		promos.DELETE("/:id", mm.RequirePermission("promo:delete"), promoHandler.RemovePromo)
		promos.GET("/types", promoHandler.ListPromoTypes)
		promos.PUT("/status", mm.RequirePermission("promo:edit"), promoHandler.SetStatusPromo)
	}
}
