package auth_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/pkg/jwt"
	"wtm-backend/pkg/logger"
)

func (au *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*authdto.LoginResponse, error) {

	dataClaim, err := jwt.ParseToken(ctx, refreshToken, au.config.RefreshSecret)
	if err != nil {
		logger.Error(ctx, "Error parsing refresh token", err.Error())
		return nil, errors.New("invalid refresh token")
	}

	if dataClaim.Type != "refresh" {
		logger.Warn(ctx,
			"Refresh token is not valid")
		return nil, errors.New("invalid refresh token")
	}

	dataUser := au.middleware.GenerateUserFromClaimToken(dataClaim)

	// Generate new JWT token
	token, err := jwt.GenerateAccessToken(dataUser, au.config.JWTSecret, au.config.DurationAccessToken)
	if err != nil {
		logger.Error(ctx, "Error generating access token", err.Error())
		return nil, err
	}

	// Set access token in Redis
	if err := au.authRepo.SetAccessToken(ctx, dataUser.ID, token, au.config.DurationAccessToken); err != nil {
		logger.Error(ctx, "Error setting access token in Redis", err.Error())
		return nil, errors.New("failed to set access token")
	}

	resp := &authdto.LoginResponse{
		Token: token,
		User: &entity.UserMin{
			ID:          dataUser.ID,
			Username:    dataUser.Username,
			Role:        dataUser.RoleName,
			Permissions: dataUser.Permissions,
			PhotoURL:    dataUser.PhotoSelfie,
			FullName:    dataUser.FullName,
		},
	}

	return resp, nil
}
