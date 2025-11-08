package auth_repository

import (
	"context"
	"time"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) CreatePasswordResetToken(ctx context.Context, userID uint, token string, expiry time.Duration) error {
	db := ar.db.GetTx(ctx)

	forgotPass := &model.PasswordResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(expiry),
	}

	if err := db.WithContext(ctx).Create(forgotPass).Error; err != nil {
		logger.Error(ctx, "Error creating password reset token:", err.Error())
		return err
	}

	return nil
}
