package promo_group_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) RemovePromoGroup(ctx context.Context, promoGroupID uint) error {
	if err := pgu.promoGroupRepo.DeletePromoGroup(ctx, promoGroupID); err != nil {
		logger.Error(ctx, "Error removing promo group:", err.Error())
		return err
	}

	return nil
}
