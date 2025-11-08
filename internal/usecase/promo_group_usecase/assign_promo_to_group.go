package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) AssignPromoToGroup(ctx context.Context, req *promogroupdto.AssignPromoToGroupRequest) error {

	if err := pgu.promoGroupRepo.AssignPromoToGroup(ctx, req.PromoGroupID, req.PromoID); err != nil {
		logger.Error(ctx, "Error assigning promo to group", err.Error())
		return err
	}

	return nil
}
