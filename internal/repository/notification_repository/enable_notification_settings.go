package notification_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (nr *NotificationRepository) EnableNotificationSettings(ctx context.Context, userID uint, channel string, types []string) error {
	db := nr.db.GetTx(ctx)

	if err := db.WithContext(ctx).
		Model(&model.UserNotificationSetting{}).
		Where("user_id = ? AND channel = ? AND type IN ?", userID, channel, types).
		Update("is_enabled", true).Error; err != nil {
		logger.Error(ctx, "Error enabling notification settings", err.Error())
		return err
	}

	return nil
}
