package notification_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (nr *NotificationRepository) GetNotificationsByUserID(ctx context.Context, filter filter.NotifFilter) ([]entity.Notification, int64, error) {
	db := nr.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", filter.UserID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "failed to count notifications", err.Error())
		return nil, 0, err
	}

	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}
	query = query.Order("is_read ASC").Order("created_at DESC")

	var notifications []model.Notification
	if err := query.Find(&notifications).Error; err != nil {
		logger.Error(ctx, "failed to get notifications by user Id", err.Error())
		return nil, 0, err
	}

	var result []entity.Notification
	if err := utils.CopyStrict(&result, &notifications); err != nil {
		logger.Error(ctx, "failed to copy notification models to entities", err.Error())
		return nil, 0, err
	}

	return result, total, nil
}
