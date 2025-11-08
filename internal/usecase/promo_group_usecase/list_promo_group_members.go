package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) ListPromoGroupMembers(ctx context.Context, req *promogroupdto.ListPromoGroupMemberRequest) (*promogroupdto.ListPromoGroupMemberResponse, int64, error) {
	promoGroupMembers, total, err := pgu.promoGroupRepo.GetPromoGroupMembers(ctx, req.PromoGroupID, req.Limit, req.Page)
	if err != nil {
		logger.Error(ctx, "Error getting promo group members", err.Error())
		return nil, total, err
	}

	respData := make([]promogroupdto.ListPromoGroupMemberData, 0, len(promoGroupMembers))
	for _, member := range promoGroupMembers {
		respData = append(respData, promogroupdto.ListPromoGroupMemberData{
			ID:           member.ID,
			Name:         member.FullName,
			AgentCompany: member.AgentCompanyName,
		})
	}

	response := &promogroupdto.ListPromoGroupMemberResponse{
		PromoGroupMembers: respData,
	}

	return response, total, nil
}
