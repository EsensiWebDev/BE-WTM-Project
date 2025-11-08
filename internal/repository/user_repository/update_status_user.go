package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) UpdateStatusUser(ctx context.Context, id uint, status uint) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Debug().
		Model(&model.User{}).
		Where("id = ?", id).
		Update("status_id", status).Error
	if err != nil {
		logger.Error(ctx, "Error updating user status:", err.Error())
		return nil, err
	}

	var user model.User
	if err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		logger.Error(ctx, "Failed to get user by Id after update", err.Error())
		return nil, err
	}

	var result entity.User
	if err := utils.CopyStrict(&result, &user); err != nil {
		logger.Error(ctx, "Failed to copy user model to entity", err.Error())
		return nil, err
	}

	return &result, nil

}
