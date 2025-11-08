package auth_repository

import (
	"context"
	"strconv"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) ValidateAccessToken(ctx context.Context, userID uint, accessToken string) (bool, error) {
	token, err := ar.redisClient.Get(ctx, "access_token:"+strconv.Itoa(int(userID)))
	if err != nil {
		logger.Error(ctx, "Error checking access token existence", "userID", userID, "err", err.Error())
		return false, err
	}
	if token != accessToken {
		logger.Warn(ctx,
			"Access token does not match", "userID", userID, "expected", accessToken, "found", token)
		return false, nil
	}
	return true, nil
}
