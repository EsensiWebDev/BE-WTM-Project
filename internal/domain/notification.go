package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/repository/filter"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *entity.Notification) error
	ReadNotification(ctx context.Context, ids []uint) error
	GetNotificationsByUserID(ctx context.Context, filter filter.NotifFilter) ([]entity.Notification, int64, error)
	DisableNotificationSettingsByChannel(ctx context.Context, userID uint, channel string) error
	EnableNotificationSettings(ctx context.Context, userID uint, channel string, types []string) error
}

type NotificationUsecase interface {
	ListNotifications(ctx context.Context, req *notifdto.ListNotificationsRequest) (*notifdto.ListNotificationsResponse, error)
	ReadNotification(ctx context.Context, req *notifdto.ReadNotificationRequest) error
	UpdateNotificationSetting(ctx context.Context, req *notifdto.UpdateNotificationSettingRequest) error
}
