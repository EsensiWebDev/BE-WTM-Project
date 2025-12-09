package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var user model.User
	if err := db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error; err != nil {
		if ur.db.ErrRecordNotFound(ctx, err) {
			logger.Info(ctx, "User not found with phone:", phone)
			return nil, nil
		}
		logger.Error(ctx, "Error getting user by email:", err.Error())
		return nil, err
	}

	var userEntity entity.User
	if err := utils.CopyPatch(&userEntity, &user); err != nil {
		logger.Error(ctx, "Error mapping user model to entity:", err.Error())
		return nil, err
	}

	return &userEntity, nil
}
