package notification_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (nr *NotificationRepository) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	db := nr.db.GetTx(ctx)

	var notificationModel model.Notification
	if err := utils.CopyStrict(&notificationModel, notification); err != nil {
		logger.Error(ctx, "Failed to copy notification entity to model", err.Error())
		return err
	}

	if err := db.Create(&notificationModel).Error; err != nil {
		logger.Error(ctx, "Failed to create notification", err.Error())
		return err
	}

	return nil
}
