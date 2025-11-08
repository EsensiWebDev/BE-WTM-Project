package promo_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (pr *PromoRepository) DeletePromo(ctx context.Context, promoID uint) error {
	db := pr.db.GetTx(ctx)

	if err := db.WithContext(ctx).Delete(&model.Promo{}, promoID).Error; err != nil {
		logger.Error(ctx, "Error deleting promo", err.Error())
		return err
	}

	return nil
}
