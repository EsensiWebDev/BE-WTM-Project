package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var modelUser model.User
	if err := utils.CopyStrict(&modelUser, user); err != nil {
		logger.Error(ctx, "Error copying user entity to model", err.Error())
		return nil, err
	}

	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(modelUser).Error
	if err != nil {
		logger.Error(ctx, "Error to update user", err.Error())
		return nil, err
	}

	return user, nil
}
