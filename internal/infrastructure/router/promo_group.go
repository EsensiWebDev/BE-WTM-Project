package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/promo_group_handler"

	"github.com/gin-gonic/gin"
)

func PromoGroupRoutes(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	promoGroupHandler := promo_group_handler.NewPromoGroupHandler(app.Usecases.PromoGroupUsecase)

	promoGroups := routerGroup.Group("/promo-groups", mm.Auth)
	{
		promoGroups.GET("", mm.RequirePermission("promo:view"), promoGroupHandler.ListPromoGroups)
		promoGroups.POST("", mm.RequirePermission("promo:create"), promoGroupHandler.CreatePromoGroup)
		promoGroups.GET("/:id", mm.RequirePermission("promo:view"), promoGroupHandler.DetailPromoGroup)
		promoGroups.DELETE("/:id", mm.RequirePermission("promo:delete"), promoGroupHandler.RemovePromoGroup)
		promoGroups.GET("/members", mm.RequirePermission("promo:view"), promoGroupHandler.ListPromoGroupMembers)
		promoGroups.POST("/members", mm.RequirePermission("promo:edit"), promoGroupHandler.AssignPromoGroupMember)
		promoGroups.DELETE("/members", mm.RequirePermission("promo:edit"), promoGroupHandler.RemovePromoGroupMember)
		promoGroups.GET("/unassigned-promos", mm.RequirePermission("promo:view"), promoGroupHandler.ListUnassignedPromos)
		promoGroups.GET("/promos", mm.RequirePermission("promo:view"), promoGroupHandler.ListPromoGroupPromos)
		promoGroups.POST("/promo", mm.RequirePermission("promo:edit"), promoGroupHandler.AssignPromoToGroup)
		promoGroups.DELETE("/promo", mm.RequirePermission("promo:edit"), promoGroupHandler.RemovePromoFromGroup)
	}
}
