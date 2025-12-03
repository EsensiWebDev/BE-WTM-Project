package promo_group_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) CreatePromoGroup(ctx context.Context, name string) error {
	promoGroup := &entity.PromoGroup{
		Name: name,
	}

	if pgu.promoGroupRepo.CheckPromoGroupExists(ctx, name) {
		logger.Error(ctx, "Promo group already exists", name)
		return errors.New("promo group already exists")
	}

	err := pgu.promoGroupRepo.CreatePromoGroup(ctx, promoGroup)
	if err != nil {
		logger.Error(ctx, "Error when creating promo group", err.Error())
		return err
	}

	return nil
}
