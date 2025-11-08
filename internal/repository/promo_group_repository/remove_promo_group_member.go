package promo_group_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (pgr *PromoGroupRepository) RemovePromoGroupMember(ctx context.Context, promoGroupID uint, memberID uint) error {
	db := pgr.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND promo_group_id = ?", memberID, promoGroupID).
		Update("promo_group_id", nil).Error
	if err != nil {
		logger.Error(ctx, "Error removing promo group member", err.Error())
		return err
	}

	return nil
}
