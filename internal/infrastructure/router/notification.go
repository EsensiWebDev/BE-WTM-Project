package router

import (
	"wtm-backend/internal/bootstrap"
	"wtm-backend/internal/handler/notification_handler"

	"github.com/gin-gonic/gin"
)

func NotificationRouter(app *bootstrap.Application, middlewareMap MiddlewareMap, routerGroup *gin.RouterGroup) {
	notificationHandler := notification_handler.NewNotificationHandler(app.Usecases.NotificationUsecase)

	notificationRouter := routerGroup.Group("/notifications", middlewareMap.TimeoutFast, middlewareMap.Auth)
	{
		notificationRouter.GET("", notificationHandler.ListNotifications)
		notificationRouter.PUT("/settings", notificationHandler.UpdateNotificationSettings)
	}
}
