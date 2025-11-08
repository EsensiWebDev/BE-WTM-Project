package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) RemovePromoFromGroup(ctx context.Context, req *promogroupdto.RemovePromoFromGroupRequest) error {

	if err := pgu.promoGroupRepo.RemovePromoFromGroup(ctx, req.PromoGroupID, req.PromoID); err != nil {
		logger.Error(ctx, "Error removing promo from group", err.Error())
		return err
	}

	return nil
}
