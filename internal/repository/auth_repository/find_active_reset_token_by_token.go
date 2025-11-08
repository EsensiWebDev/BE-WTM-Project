package auth_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ar *AuthRepository) FindActiveResetTokenByToken(ctx context.Context, token string) (string, error) {
	db := ar.db.GetTx(ctx)

	var user model.User
	var email string

	err := db.WithContext(ctx).
		Table("password_reset_tokens AS prt").
		Select("u.email").
		Joins("JOIN users u ON u.id = prt.user_id").
		Where("prt.token = ?", token).
		Where("prt.expires_at > NOW()").
		Where("prt.used = FALSE").
		First(&user).Error

	if err != nil {
		logger.Error(ctx, "Error finding user by reset token:", err.Error())
		return email, err
	}

	return user.Email, nil
}
