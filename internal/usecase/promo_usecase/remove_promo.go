package promo_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) RemovePromo(ctx context.Context, promoID string) error {
	promo, err := pu.promoRepo.GetPromoByExternalID(ctx, promoID)
	if err != nil {
		logger.Error(ctx, "Error getting promo by Id", "error", err, "promoID", promoID)
		return err
	}
	if err := pu.promoRepo.DeletePromo(ctx, promo.ID); err != nil {
		logger.Error(ctx, "Error removing promo", "error", err, "promoID", promoID)
		return err
	}

	return nil
}
