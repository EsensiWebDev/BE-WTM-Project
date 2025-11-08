package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListAllBedTypes(ctx context.Context, req *hoteldto.ListAllBedTypesRequest) (*hoteldto.ListAllBedTypesResponse, error) {
	filterRepo := &filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	bedTypes, total, err := hu.hotelRepo.GetBedTypes(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting bed types", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListAllBedTypesResponse{
		BedTypes: bedTypes,
		Total:    total,
	}

	return resp, nil
}
