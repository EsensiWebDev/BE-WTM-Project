package router

import (
	"github.com/gin-gonic/gin"
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/email_handler"
)

func EmailRoute(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	emailHandler := email_handler.NewEmailHandler(app.Usecases.EmailUsecase)

	emailRouter := routerGroup.Group("/email", middlewareMap.TimeoutFast)
	{
		emailRouter.POST("/contact-us", emailHandler.SendContactUs)
		emailRouter.GET("/template", middlewareMap.Auth, emailHandler.EmailTemplate)
		emailRouter.PUT("/template", middlewareMap.Auth, emailHandler.UpdateEmailTemplate)
		emailRouter.GET("/logs", middlewareMap.Auth, emailHandler.ListEmailLogs)
	}

}
