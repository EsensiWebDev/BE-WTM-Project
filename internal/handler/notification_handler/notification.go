package notification_handler

import "wtm-backend/internal/domain"

type NotificationHandler struct {
	notifUsecase domain.NotificationUsecase
}

func NewNotificationHandler(notifUsecase domain.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{notifUsecase: notifUsecase}
}
