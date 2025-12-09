package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/email_handler"

	"github.com/gin-gonic/gin"
)

func EmailRoute(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	emailHandler := email_handler.NewEmailHandler(app.Usecases.EmailUsecase)

	emailRouter := routerGroup.Group("/email")
	{
		emailRouter.POST("/contact-us", middlewareMap.TimeoutFast, emailHandler.SendContactUs)
		emailRouter.GET("/template", middlewareMap.Auth, middlewareMap.TimeoutFast, emailHandler.EmailTemplate)
		emailRouter.PUT("/template", middlewareMap.Auth, middlewareMap.TimeoutFile, emailHandler.UpdateEmailTemplate)
		emailRouter.GET("/logs", middlewareMap.Auth, middlewareMap.TimeoutFast, emailHandler.ListEmailLogs)
	}

}
