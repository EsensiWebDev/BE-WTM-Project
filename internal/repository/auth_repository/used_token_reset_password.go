package auth_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) UsedTokenResetPassword(ctx context.Context, token string) (uint, error) {
	db := ar.db.GetTx(ctx)

	var trp model.PasswordResetToken
	// Ambil semua kolom supaya ID ikut terisi
	if err := db.Where("token = ? AND used = FALSE", token).First(&trp).Error; err != nil {
		if ar.db.ErrRecordNotFound(ctx, err) {
			logger.Error(ctx, "Password reset token not found", "token", token, "err", err.Error())
			return 0, errors.New("invalid or expired token")
		}
		logger.Error(ctx, "Error finding password reset token", "token", token, "err", err.Error())
		return 0, err
	}

	// Update langsung kolom used = true
	if err := db.Model(&trp).Update("used", true).Error; err != nil {
		logger.Error(ctx, "Error updating password reset token as used", "token", token, "err", err.Error())
		return 0, err
	}

	return trp.UserID, nil
}
