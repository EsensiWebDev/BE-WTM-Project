package promo_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) PromoByID(ctx context.Context, promoID string) (*entity.PromoWithExternalID, error) {
	promo, err := pu.promoRepo.GetPromoByExternalID(ctx, promoID)
	if err != nil {
		logger.Error(ctx, "Error getting promo by Id", "error", err, "promoID", promoID)
		return nil, err
	}

	promoEntity, err := pu.promoRepo.GetPromoByID(ctx, promo.ID, nil)
	if err != nil {
		logger.Error(ctx, "Error getting promo by Id", "error", err, "promoID", promoID)
		return nil, err
	}

	resp := &entity.PromoWithExternalID{
		ID:             promoEntity.ExternalID,
		Name:           promoEntity.Name,
		StartDate:      promoEntity.StartDate,
		EndDate:        promoEntity.EndDate,
		Code:           promoEntity.Code,
		Description:    promoEntity.Description,
		PromoTypeID:    promoEntity.PromoTypeID,
		Detail:         promoEntity.Detail,
		IsActive:       promoEntity.IsActive,
		PromoTypeName:  promoEntity.PromoTypeName,
		PromoGroups:    promoEntity.PromoGroups,
		PromoRoomTypes: promoEntity.PromoRoomTypes,
		PromoGroupIDs:  promoEntity.PromoGroupIDs,
		Duration:       promoEntity.Duration,
	}

	return resp, nil
}
