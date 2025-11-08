package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var user model.User
	err := db.WithContext(ctx).
		Preload("Status").
		Preload("Role").
		Preload("Role.Permissions").
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if ur.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "User not found with username", username)
			return nil, nil
		}
		logger.Error(ctx, "Error to get user by username", err.Error())
		return nil, err
	}

	var entityUser entity.User
	if err := utils.CopyPatch(&entityUser, user); err != nil {
		logger.Error(ctx, "Error copying user model to entity", err.Error())
		return nil, err
	}

	entityUser.StatusName = user.Status.Status

	if user.Role != nil {
		entityUser.RoleName = user.Role.Role
		if user.StatusID == constant.StatusUserActiveID {
			for _, permission := range user.Role.Permissions {
				entityUser.Permissions = append(entityUser.Permissions, permission.Permission)
			}
		}
	}

	return &entityUser, nil
}
