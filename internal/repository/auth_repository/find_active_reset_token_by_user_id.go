package auth_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) FindActiveResetTokenByUserID(ctx context.Context, userID uint) (string, error) {
	db := ar.db.GetTx(ctx)

	var prt model.PasswordResetToken
	err := db.WithContext(ctx).
		Select("token").
		Where("user_id = ?", userID).
		Where("expires_at > NOW()").
		Where("used = FALSE").
		First(&prt).Error
	if err != nil {
		logger.Error(ctx, "Error finding active reset token by user ID:", err.Error())
		return "", err
	}

	return prt.Token, nil
}
