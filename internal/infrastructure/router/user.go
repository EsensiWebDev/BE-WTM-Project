package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/user_handler"
	"wtm-backend/pkg/constant"
)

func UserRoutes(app *bootstrap.Application, mm MiddlewareMap, routerGroup *gin.RouterGroup) {
	userHandler := user_handler.NewUserHandler(app.Usecases.UserUsecase, app.Config)

	routerGroup.POST("/register", mm.TimeoutFile, userHandler.Register)

	roleAccess := routerGroup.Group("/role-access", mm.Auth, mm.RequireRole(constant.RoleSuperAdminCap))
	{
		roleAccess.GET("", userHandler.ListRoleAccess)
		roleAccess.PUT("", userHandler.UpdateRoleAccess)
	}

	profile := routerGroup.Group("/profile", mm.Auth)
	{
		profile.GET("", userHandler.Profile)
		profile.PUT("", mm.TimeoutSlow, userHandler.UpdateProfile)
		profile.PUT("/setting", userHandler.UpdateSetting)
		profile.PUT("/file", userHandler.UpdateFile)
	}

	users := routerGroup.Group("/users", mm.Auth)
	{
		users.GET("/control", mm.RequirePermission("account:view"), userHandler.ListControlUsers)
		users.GET("", mm.RequirePermission("account:view"), userHandler.ListUsers)
		users.PUT("", mm.RequirePermission("account:edit"), userHandler.UpdateUserByAdmin)
		users.POST("", mm.RequirePermission("account:create"), userHandler.CreateUserByAdmin)
		users.GET("/agent-companies", userHandler.ListAgentCompanies)
		users.GET("/by-agent-company/:id", userHandler.ListUsersByAgentCompany)
		users.GET("/status", userHandler.ListStatusUsers)
		users.POST("/status", mm.RequirePermission("account:edit"), userHandler.UpdateStatusUser)
	}

}
