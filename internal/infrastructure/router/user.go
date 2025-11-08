package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/user_handler"
)

func UserRoutes(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	userHandler := user_handler.NewUserHandler(app.Usecases.UserUsecase, app.Config)

	routerGroup.POST("/register", middlewareMap.TimeoutFile, userHandler.Register)

	roleAccess := routerGroup.Group("/role-access", middlewareMap.Auth)
	{
		roleAccess.GET("", userHandler.ListRoleAccess)
		roleAccess.PUT("", userHandler.UpdateRoleAccess)
	}

	profile := routerGroup.Group("/profile", middlewareMap.Auth)
	{
		profile.GET("", userHandler.Profile)
		profile.PUT("", middlewareMap.TimeoutSlow, userHandler.UpdateProfile)
		profile.PUT("/setting", userHandler.UpdateSetting)
		profile.PUT("/file", userHandler.UpdateFile)
	}

	users := routerGroup.Group("/users", middlewareMap.Auth)
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/control", userHandler.ListControlUsers)
		users.PUT("", userHandler.UpdateUserByAdmin)
		users.POST("", userHandler.CreateUserByAdmin)
		users.GET("/agent-companies", userHandler.ListAgentCompanies)
		users.GET("/by-agent-company/:id", userHandler.ListUsersByAgentCompany)
		//users.GET("/by-role/:role", userHandler.ListUsersByRole)
		users.GET("/status", userHandler.ListStatusUsers)
		users.POST("/status", userHandler.UpdateStatusUser)
	}

}
