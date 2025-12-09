package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ur *UserRepository) GetNotificationSettings(ctx context.Context, userID uint) ([]entity.UserNotificationSetting, error) {
	db := ur.db.GetTx(ctx)

	var notificationSettings []model.UserNotificationSetting
	err := db.Where("user_id = ?", userID).Find(&notificationSettings).Error
	if err != nil {
		logger.Error(ctx, "GetNotificationSettings", err.Error())
		return nil, err
	}

	var notificationSettingEntities []entity.UserNotificationSetting
	notificationSettingEntities = make([]entity.UserNotificationSetting, 0, len(notificationSettings))
	for _, setting := range notificationSettings {
		notificationSettingEntities = append(notificationSettingEntities, entity.UserNotificationSetting{
			UserID:    setting.UserID,
			Channel:   setting.Channel,
			Type:      setting.Type,
			IsEnabled: setting.IsEnabled,
		})
	}

	return notificationSettingEntities, nil
}
