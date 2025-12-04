package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/auth_handler"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	authHandler := auth_handler.NewAuthHandler(app.Usecases.AuthUsecase, app.Config)

	auth := routerGroup.Group("")
	{
		auth.POST("/login", authHandler.Login)
		auth.GET("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", middlewareMap.Auth, authHandler.Logout)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.GET("/reset-password", authHandler.ValidateTokenResetPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

}
