package promo_usecase

import (
	"context"
	"time"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/repository/filter"
)

func (pu *PromoUsecase) ListPromos(ctx context.Context, req *promodto.ListPromosRequest) (*promodto.ListPromosResponse, error) {
	filterReq := &filter.DefaultFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	promoEntity, total, err := pu.promoRepo.GetPromos(ctx, filterReq)
	if err != nil {
		return nil, err
	}

	var promos []promodto.PromoResponse
	for _, promo := range promoEntity {

		promos = append(promos, promodto.PromoResponse{
			ID:               promo.ID,
			PromoName:        promo.Name,
			PromoCode:        promo.Code,
			Duration:         promo.PromoRoomTypes[0].TotalNights,
			PromoStartDate:   promo.StartDate.Format(time.RFC3339),
			PromoEndDate:     promo.EndDate.Format(time.RFC3339),
			IsActive:         promo.IsActive,
			PromoType:        promo.PromoTypeName,
			PromoDetail:      promo.Detail,
			PromoDescription: promo.Description,
		})
	}

	resp := &promodto.ListPromosResponse{
		Promos: promos,
		Total:  total,
	}

	return resp, nil
}
