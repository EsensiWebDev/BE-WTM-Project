package promo_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (pr *PromoRepository) UpdatePromoStatus(ctx context.Context, promoID uint, isActive bool) error {
	db := pr.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Model(&model.Promo{}).
		Where("id = ?", promoID).
		Update("is_active", isActive).Error
	if err != nil {
		logger.Error(ctx, "Error setting promo status", err.Error())
		return err
	}

	return nil
}
