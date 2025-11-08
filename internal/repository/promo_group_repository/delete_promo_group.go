package promo_group_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (pgr *PromoGroupRepository) DeletePromoGroup(ctx context.Context, promoGroupID uint) error {
	db := pgr.db.GetTx(ctx)

	if err := db.WithContext(ctx).Delete(&model.PromoGroup{}, promoGroupID).Error; err != nil {
		logger.Error(ctx, "Error deleting promo", err.Error())
		return err
	}

	return nil
}
