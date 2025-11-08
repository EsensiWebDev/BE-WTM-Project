package auth_repository

import (
	"context"
	"strings"
	"time"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) SetEmailForgotPassword(ctx context.Context, email string, duration time.Duration) error {
	key := "forgot_password:" + strings.ToLower(email)
	err := ar.redisClient.Set(ctx, key, "1", duration)
	if err != nil {
		logger.Error(ctx, "Error setting email forgot password in redis:", err.Error())
		return err
	}
	return nil
}
