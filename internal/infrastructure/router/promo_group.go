package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/promo_group_handler"
)

func PromoGroupRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	promoGroupHandler := promo_group_handler.NewPromoGroupHandler(app.Usecases.PromoGroupUsecase)

	promoGroups := routerGroup.Group("/promo-groups", middlewareMap.Auth)
	{
		promoGroups.GET("", promoGroupHandler.ListPromoGroups)
		promoGroups.POST("", promoGroupHandler.CreatePromoGroup)
		promoGroups.GET("/:id", promoGroupHandler.DetailPromoGroup)
		promoGroups.DELETE("/:id", promoGroupHandler.RemovePromoGroup)
		promoGroups.GET("/members", promoGroupHandler.ListPromoGroupMembers)
		promoGroups.POST("/members", promoGroupHandler.AssignPromoGroupMember)
		promoGroups.DELETE("/members", promoGroupHandler.RemovePromoGroupMember)
		promoGroups.GET("/unassigned-promos", promoGroupHandler.ListUnassignedPromos)
		promoGroups.GET("/promos", promoGroupHandler.ListPromoGroupPromos)
		promoGroups.POST("/promo", promoGroupHandler.AssignPromoToGroup)
		promoGroups.DELETE("/promo", promoGroupHandler.RemovePromoFromGroup)
	}
}
