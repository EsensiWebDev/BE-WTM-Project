package user_usecase

import (
	"context"
	"errors"
	"strings"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (uu *UserUsecase) UpdateSetting(ctx context.Context, req *userdto.UpdateSettingRequest) error {

	dataUser, err := uu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error generating user from context", err.Error())
		return err
	}
	if dataUser == nil {
		logger.Error(ctx, "User not found in context")
		return errors.New("user not found in context")
	}

	logger.Info(ctx, "Update setting user request", req)
	logger.Info(ctx, "Update setting user dataUser", dataUser)

	if dataUser.Username != req.Username {
		user, err := uu.userRepo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			logger.Error(ctx, "Error getting user by username", err.Error())
			return err
		}

		if user != nil {
			logger.Error(ctx, "Username already exists", req.Username)
			return errors.New("username already exists")
		}
	} else {
		if strings.TrimSpace(req.NewPassword) == "" {
			logger.Error(ctx, "New password cannot be empty")
			return errors.New("new password cannot be empty")
		}
		if req.NewPassword == req.OldPassword {
			logger.Error(ctx, "New password cannot be the same as old password")
			return errors.New("new password cannot be the same as old password")
		}
	}

	userDB, err := uu.userRepo.GetUserByID(ctx, dataUser.ID)
	if err != nil {
		logger.Error(ctx, "Error getting user by Id", err.Error())
		return err
	}

	if userDB == nil {
		logger.Error(ctx, "User not found in database")
		return errors.New("user not found in database")
	}

	// Verifikasi password
	if !utils.ComparePassword(ctx, userDB.Password, req.OldPassword) {
		logger.Warn(ctx,
			"Password is invalid")
		return errors.New("password is invalid")
	}

	if strings.TrimSpace(req.NewPassword) != "" {
		passCrypt, err := utils.GeneratePassword(ctx, req.NewPassword)
		if err != nil {
			logger.Error(ctx, "Error generating password", err.Error())
			return err
		}

		userDB.Password = passCrypt
	}

	userDB.Username = req.Username
	_, err = uu.userRepo.UpdateUser(ctx, userDB)
	if err != nil {
		logger.Error(ctx, "Error updating user password", err.Error())
		return err
	}

	if err := uu.authRepo.DeleteAccessToken(ctx, dataUser.ID); err != nil {
		logger.Error(ctx, "Error to delete access token", err.Error())
		return errors.New("failed to delete access token")
	}

	return nil

}
