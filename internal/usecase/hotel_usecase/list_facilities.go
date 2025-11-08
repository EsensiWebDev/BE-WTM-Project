package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListFacilities(ctx context.Context, req *hoteldto.ListFacilitiesRequest) (*hoteldto.ListFacilitiesResponse, error) {

	filterRepo := &filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	facilities, total, err := hu.hotelRepo.GetFacilities(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting facilities", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListFacilitiesResponse{
		Facilities: facilities,
		Total:      total,
	}

	return resp, nil
}
