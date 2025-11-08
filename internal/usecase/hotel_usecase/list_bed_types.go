package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListBedTypes(ctx context.Context, roomTypeID uint) (*hoteldto.ListBedTypeResponse, error) {
	bedTypes, err := hu.hotelRepo.GetBedTypeByRoomTypeID(ctx, roomTypeID)
	if err != nil {
		logger.Error(ctx, "Error getting bed types by room type", "roomTypeID", roomTypeID, "err", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListBedTypeResponse{
		BedTypes: bedTypes,
	}

	return resp, nil
}
