package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) RemovePromoGroupMember(ctx context.Context, req *promogroupdto.RemovePromoGroupMemberRequest) error {
	err := pgu.promoGroupRepo.RemovePromoGroupMember(ctx, req.PromoGroupID, req.MemberID)
	if err != nil {
		logger.Error(ctx, "Error removing promo group member", err.Error())
		return err
	}
	return nil
}
