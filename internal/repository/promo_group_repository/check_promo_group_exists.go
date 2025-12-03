package promo_group_repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (pgr *PromoGroupRepository) CheckPromoGroupExists(ctx context.Context, name string) bool {
	db := pgr.db.GetTx(ctx)

	var promoGroup model.PromoGroup
	err := db.WithContext(ctx).
		Where("LOWER(name) = ?", strings.TrimSpace(strings.ToLower(name))).
		Debug().
		First(&promoGroup).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	if err != nil {
		logger.Error(ctx, "Error checking promo group exists", err.Error())
		return false
	}
	return true
}
