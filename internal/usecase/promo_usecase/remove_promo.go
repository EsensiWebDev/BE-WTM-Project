package promo_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) RemovePromo(ctx context.Context, promoID uint) error {
	if err := pu.promoRepo.DeletePromo(ctx, promoID); err != nil {
		logger.Error(ctx, "Error removing promo", "error", err, "promoID", promoID)
		return err
	}

	return nil
}
