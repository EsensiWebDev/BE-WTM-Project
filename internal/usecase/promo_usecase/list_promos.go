package promo_usecase

import (
	"context"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/repository/filter"
)

func (pu *PromoUsecase) ListPromos(ctx context.Context, req *promodto.ListPromosRequest) (*promodto.ListPromosResponse, int64, error) {
	filterReq := &filter.DefaultFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	promoEntity, total, err := pu.promoRepo.GetPromos(ctx, filterReq)
	if err != nil {
		return nil, 0, err
	}

	var promos []promodto.PromoResponse
	for _, promo := range promoEntity {

		promos = append(promos, promodto.PromoResponse{
			ID:               promo.ID,
			PromoName:        promo.Name,
			PromoCode:        promo.Code,
			Duration:         promo.Duration,
			PromoStartDate:   promo.StartDate.String(),
			PromoEndDate:     promo.EndDate.String(),
			IsActive:         promo.IsActive,
			PromoType:        promo.PromoTypeName,
			PromoDetail:      promo.Detail,
			PromoDescription: promo.Description,
		})
	}

	resp := &promodto.ListPromosResponse{
		Promos: promos,
	}

	return resp, total, nil
}
