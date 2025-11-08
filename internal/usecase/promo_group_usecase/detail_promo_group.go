package promo_group_usecase

import (
	"context"
	"strconv"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) DetailPromoGroup(ctx context.Context, promoGroupID uint) (*entity.PromoGroup, error) {
	promoGroup, err := pgu.promoGroupRepo.GetPromoGroupByID(ctx, promoGroupID)
	if err != nil {
		logger.Error(ctx, "Error getting promo group by ID:", err.Error())
		return nil, err
	}
	if promoGroup == nil {
		logger.Error(ctx, "Promo group not found with ID:", strconv.Itoa(int(promoGroupID)))
		return nil, nil
	}

	return promoGroup, nil
}
