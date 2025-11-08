package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hu *HotelUsecase) ListProvinces(ctx context.Context, req *hoteldto.ListProvincesRequest) (*hoteldto.ListProvincesResponse, error) {
	filterRepo := filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	provinces, total, err := hu.hotelRepo.GetProvinces(ctx, &filterRepo)
	if err != nil {
		logger.Error(ctx, "Failed to get provinces", err.Error())
		return nil, err
	}

	var provincesCapitalized []string
	for _, province := range provinces {
		prov := utils.CapitalizeWords(province)
		provincesCapitalized = append(provincesCapitalized, prov)
	}

	resp := &hoteldto.ListProvincesResponse{
		Provinces: provincesCapitalized,
		Total:     total,
	}

	return resp, nil
}
