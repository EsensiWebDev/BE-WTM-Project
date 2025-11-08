package auth_repository

import (
	"context"
	"strings"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) DeleteEmailForgotPassword(ctx context.Context, email string) error {
	key := "forgot_password:" + strings.ToLower(email)
	if err := ar.redisClient.Delete(ctx, key); err != nil {
		logger.Error(ctx, "Error deleting forgot password email from redis", "email", email, "err", err.Error())
		return err
	}

	return nil
}
