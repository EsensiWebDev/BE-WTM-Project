package auth_repository

import (
	"context"
	"strconv"
	"time"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) SetAccessToken(ctx context.Context, userID uint, accessToken string, expiry time.Duration) error {
	// set access token to redis with key "access_token:<user_id>"
	err := ar.redisClient.Set(ctx, "access_token:"+strconv.Itoa(int(userID)), accessToken, expiry)
	if err != nil {
		logger.Error(ctx, "Error setting access token in redis", "userID", userID, "err", err.Error())
		return err
	}
	return nil
}
