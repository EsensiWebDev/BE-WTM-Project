package notification_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (nr *NotificationRepository) ReadNotification(ctx context.Context, ids []uint) error {
	db := nr.db.GetTx(ctx)

	if err := db.Model(&entity.Notification{}).
		Where("id IN ?", ids).
		Update("is_read", true).Error; err != nil {
		logger.Error(ctx, "Failed to mark notification as read", err.Error())
		return err
	}

	return nil
}
