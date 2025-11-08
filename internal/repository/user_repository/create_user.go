package user_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var modelUser model.User
	if err := utils.CopyStrict(&modelUser, user); err != nil {
		logger.Error(ctx, "Error copying user entity to model", err.Error())
		return nil, err
	}

	// Cek apakah role-nya agent
	if modelUser.RoleID == constant.RoleAgentID { // misalnya agentRoleID := 2
		defaultSettings := []model.UserNotificationSetting{
			{Channel: "email", Type: "booking", IsEnabled: true},
			{Channel: "email", Type: "reject", IsEnabled: true},
			{Channel: "web", Type: "booking", IsEnabled: true},
			{Channel: "web", Type: "reject", IsEnabled: true},
		}
		modelUser.UserNotificationSettings = defaultSettings
	}

	err := db.WithContext(ctx).Debug().Create(&modelUser).Error
	if err != nil {
		if ur.db.ErrDuplicateKey(ctx, err) {
			logger.Warn(ctx, "User already exists with username", user.Username)
			return nil, errors.New("user already exists with username: " + user.Username)
		}
		logger.Error(ctx, "Error to add user", err.Error())
		return nil, err
	}

	var entityUser entity.User
	if err := utils.CopyStrict(&entityUser, modelUser); err != nil {
		logger.Error(ctx, "Error copying user model to entity", err.Error())
		return nil, err
	}

	return &entityUser, nil
}
