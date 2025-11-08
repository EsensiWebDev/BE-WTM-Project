package auth_repository

import (
	"context"
	"strings"
	"time"
)

func (ar *AuthRepository) ValidateEmailForgotPassword(ctx context.Context, email string) (bool, time.Duration, error) {
	key := "forgot_password:" + strings.ToLower(email)

	ttl, err := ar.redisClient.TTL(ctx, key)
	if err != nil {
		return false, 0, err
	}

	// TTL bisa -2 (key tidak ada), -1 (key tidak punya expiry), atau positif
	if ttl <= 0 {
		return false, 0, nil
	}

	return true, ttl, nil
}
