package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) ListPromoGroups(ctx context.Context, req *promogroupdto.ListPromoGroupRequest) (*promogroupdto.ListPromoGroupResponse, int64, error) {
	promoGroups, total, err := pgu.promoGroupRepo.GetPromoGroups(ctx, req.Search, req.Limit, req.Page)
	if err != nil {
		logger.Error(ctx, "Error getting promo groups", err.Error())
		return nil, 0, err
	}

	resp := &promogroupdto.ListPromoGroupResponse{
		PromoGroups: promoGroups,
	}

	return resp, total, nil
}
