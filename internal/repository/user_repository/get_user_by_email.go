package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var user model.User
	if err := db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if ur.db.ErrRecordNotFound(ctx, err) {
			logger.Error(ctx, "User not found by email:", email)
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
