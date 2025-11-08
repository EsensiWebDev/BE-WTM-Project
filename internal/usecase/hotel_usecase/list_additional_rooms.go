package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListAdditionalRooms(ctx context.Context, req *hoteldto.ListAdditionalRoomsRequest) (*hoteldto.ListAdditionalRoomsResponse, error) {
	filterRepo := &filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	additionalRooms, total, err := hu.hotelRepo.GetAdditionalRooms(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting additional rooms", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListAdditionalRoomsResponse{
		AdditionalRooms: additionalRooms,
		Total:           total,
	}

	return resp, nil
}
