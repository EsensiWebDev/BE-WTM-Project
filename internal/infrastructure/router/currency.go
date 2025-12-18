package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/currency_handler"

	"github.com/gin-gonic/gin"
)

func CurrencyRoutes(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	currencyHandler := currency_handler.NewCurrencyHandler(app.Usecases.CurrencyUsecase)

	currencies := routerGroup.Group("/currencies", mm.Auth, mm.RequirePermission("settings:view"))
	{
		currencies.GET("", currencyHandler.GetAllCurrencies)
		currencies.GET("/active", currencyHandler.GetActiveCurrencies)
		currencies.POST("", mm.RequirePermission("settings:create"), currencyHandler.CreateCurrency)
		currencies.PUT("/:id", mm.RequirePermission("settings:edit"), currencyHandler.UpdateCurrency)
	}
}
