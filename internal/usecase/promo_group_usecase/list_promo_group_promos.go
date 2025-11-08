package promo_group_usecase

import (
	"context"
	"wtm-backend/internal/dto/promogroupdto"
)

func (pgu *PromoGroupUsecase) ListPromoGroupPromos(ctx context.Context, req *promogroupdto.ListPromoGroupPromosRequest) (*promogroupdto.ListPromoGroupPromosResponse, int64, error) {

	promos, total, err := pgu.promoGroupRepo.GetPromosByPromoGroupID(ctx, req.ID, req.Search, req.Limit, req.Page)
	if err != nil {
		return nil, total, err
	}

	var respData []promogroupdto.ListPromoGroupPromosData
	for _, promo := range promos {

		var startDate, endDate string
		if promo.StartDate != nil {
			startDate = promo.StartDate.Format("02 January 2006 15:04:05")
		}
		if promo.EndDate != nil {
			endDate = promo.EndDate.Format("02 January 2006 15:04:05")
		}

		respData = append(respData, promogroupdto.ListPromoGroupPromosData{
			PromoID:        promo.ID,
			PromoName:      promo.Name,
			PromoCode:      promo.Code,
			PromoStartDate: startDate,
			PromoEndDate:   endDate,
		})
	}

	resp := &promogroupdto.ListPromoGroupPromosResponse{
		Promos: respData,
	}

	return resp, total, nil
}
