package auth_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) UsedTokenResetPassword(ctx context.Context, token string) (uint, error) {
	db := ar.db.GetTx(ctx)

	var userId uint
	var trp model.PasswordResetToken
	if err := db.Select("user_id").Where("token = ? AND used = FALSE", token).First(&trp).Error; err != nil {
		if ar.db.ErrRecordNotFound(ctx, err) {
			logger.Error(ctx, "Password reset token not found", "token", token, "err", err.Error())
			return userId, errors.New("invalid or expired token")
		}
		logger.Error(ctx, "Error finding password reset token", "token", token, "err", err.Error())
		return userId, err
	}
	trp.Used = true
	if err := db.Save(&trp).Error; err != nil {
		logger.Error(ctx, "Error updating password reset token as used", "token", token, "err", err.Error())
		return userId, err
	}

	userId = trp.UserID

	return userId, nil
}
