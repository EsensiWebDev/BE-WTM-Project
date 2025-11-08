package notification_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (nr *NotificationRepository) ReadNotification(ctx context.Context, id int64) error {
	db := nr.db.GetTx(ctx)

	if err := db.Model(&entity.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error; err != nil {
		logger.Error(ctx, "Failed to mark notification as read", err.Error())
		return err
	}

	return nil
}
