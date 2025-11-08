package auth_usecase

import (
	"context"
	"errors"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/jwt"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (au *AuthUsecase) Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, string, error) {

	user, err := au.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Error(ctx, "Error to get user by username", err.Error())
		return nil, "", errors.New("user not found")
	}

	if user == nil {
		logger.Warn(ctx,
			"User not found")
		return nil, "", errors.New("user not found")
	}

	if user.StatusID != constant.StatusUserActiveID {
		logger.Warn(ctx,
			"User is not active")
		return nil, "", errors.New("user is not active")
	}

	// Verifikasi password
	if !utils.ComparePassword(ctx, user.Password, req.Password) {
		logger.Warn(ctx,
			"Password is invalid")
		return nil, "", errors.New("invalid password")
	}

	if user.PhotoSelfie != "" {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPublic)
		photoProfile, err := au.fileStorage.GetFile(ctx, bucketName, user.PhotoSelfie)
		if err != nil {
			logger.Error(ctx, "Error getting user profile photo", err.Error())
			return nil, "", fmt.Errorf("failed to get user profile photo: %s", err.Error())
		}

		user.PhotoSelfie = photoProfile
	}

	// Generate JWT token
	token, err := jwt.GenerateAccessToken(user, au.config.JWTSecret, au.config.DurationAccessToken)
	if err != nil {
		logger.Error(ctx, "Error generating access token", err.Error())
		return nil, "", err
	}

	// Generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken(user, au.config.RefreshSecret, au.config.DurationRefreshToken)
	if err != nil {
		logger.Error(ctx, "Error generating refresh token", err.Error())
		return nil, "", err
	}

	// Set access token in Redis
	if err := au.authRepo.SetAccessToken(ctx, user.ID, token, au.config.DurationAccessToken); err != nil {
		logger.Error(ctx, "Error setting access token in Redis", err.Error())
		return nil, "", errors.New("failed to set access token")
	}

	resp := &authdto.LoginResponse{
		Token: token,
		User: &entity.UserMin{
			ID:          user.ID,
			Username:    user.Username,
			RoleID:      user.RoleID,
			Role:        user.RoleName,
			Permissions: user.Permissions,
			PhotoURL:    user.PhotoSelfie,
			FullName:    user.FullName,
		},
	}

	return resp, refreshToken, nil
}
