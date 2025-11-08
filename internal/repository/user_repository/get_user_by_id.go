package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUserByID(ctx context.Context, userID uint) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var user model.User
	err := db.WithContext(ctx).
		Where("id = ?", userID).
		Preload("UserNotificationSettings").
		First(&user).Error

	if err != nil {
		if ur.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "User not found with Id", userID)
			return nil, nil
		}
		logger.Error(ctx, "Error to get user by Id", err.Error())
		return nil, err
	}

	var entityUser entity.User
	if err := utils.CopyPatch(&entityUser, user); err != nil {
		logger.Error(ctx, "Error copying user model to entity", err.Error())
		return nil, err
	}

	return &entityUser, nil
}
