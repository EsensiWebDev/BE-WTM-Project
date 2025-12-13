package promo_group_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) ListUnassignedPromos(ctx context.Context, req *promogroupdto.ListUnassignedPromosRequest) (*promogroupdto.ListUnassignedPromosResponse, error) {

	filterReq := &filter.PromoGroupFilter{
		PromoGroupID: req.PromoGroupID,
	}
	filterReq.PaginationRequest = req.PaginationRequest

	promos, total, err := pgu.promoGroupRepo.GetUnassignedPromos(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "Error listing unassigned promos:", err.Error())
		return nil, err
	}

	var datas []promogroupdto.ListUnassignedPromoData
	for _, promo := range promos {
		data := promogroupdto.ListUnassignedPromoData{
			ID:   promo.ID,
			Name: fmt.Sprintf("%s - %s", promo.Code, promo.Name),
		}

		datas = append(datas, data)
	}

	response := &promogroupdto.ListUnassignedPromosResponse{
		Promos: datas,
		Total:  total,
	}

	return response, nil
}
