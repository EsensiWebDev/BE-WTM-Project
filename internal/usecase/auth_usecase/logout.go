package auth_usecase

import (
	"context"
	"errors"
	"wtm-backend/pkg/logger"
)

func (au *AuthUsecase) Logout(ctx context.Context) error {
	dataUser, err := au.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error to get user from context", err.Error())
		return errors.New("failed to get user from context")
	}

	if err := au.authRepo.DeleteAccessToken(ctx, dataUser.ID); err != nil {
		logger.Error(ctx, "Error to delete access token", err.Error())
		return errors.New("failed to delete access token")
	}

	return nil
}
