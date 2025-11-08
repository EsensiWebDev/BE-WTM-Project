package promo_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) PromoByID(ctx context.Context, promoID uint) (*entity.Promo, error) {
	promoEntity, err := pu.promoRepo.GetPromoByID(ctx, promoID, nil)
	if err != nil {
		logger.Error(ctx, "Error getting promo by Id", "error", err, "promoID", promoID)
		return nil, err
	}

	return promoEntity, nil
}
