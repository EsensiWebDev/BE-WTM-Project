package auth_repository

import (
	"context"
	"strconv"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) DeleteAccessToken(ctx context.Context, userID uint) error {
	// Delete access token from redis
	err := ar.redisClient.Delete(ctx, "access_token:"+strconv.Itoa(int(userID)))
	if err != nil {
		logger.Error(ctx, "Error deleting access token from redis", "userID", userID, "err", err.Error())
		return err
	}
	return nil
}
