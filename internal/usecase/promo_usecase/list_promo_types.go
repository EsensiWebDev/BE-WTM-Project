package promo_usecase

import (
	"context"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) ListPromoTypes(ctx context.Context, req *promodto.ListPromoTypesRequest) (*promodto.ListPromoTypesResponse, int64, error) {

	filterRepo := &filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	promoTypes, total, err := pu.promoRepo.GetPromoTypes(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting facilities", err.Error())
		return nil, 0, err
	}

	resp := &promodto.ListPromoTypesResponse{
		PromoTypes: promoTypes,
	}

	return resp, total, nil
}
