package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
)

func (pgu *PromoGroupUsecase) CreatePromoGroup(ctx context.Context, name string) error {
	promoGroup := &entity.PromoGroup{
		Name: name,
	}
	err := pgu.promoGroupRepo.CreatePromoGroup(ctx, promoGroup)
	if err != nil {
		return err
	}

	return nil
}
