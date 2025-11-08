package notification_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (nr *NotificationRepository) DisableNotificationSettingsByChannel(ctx context.Context, userID uint, channel string) error {
	db := nr.db.GetTx(ctx)

	if err := db.WithContext(ctx).
		Model(&model.UserNotificationSetting{}).
		Where("user_id = ? AND channel = ?", userID, channel).
		Update("is_enabled", false).Error; err != nil {
		logger.Error(ctx, "failed to disable notification settings by channel", err.Error())
		return err
	}

	return nil
}
