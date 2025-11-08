package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListRoomTypes(ctx context.Context, hotelID uint) (*hoteldto.ListRoomTypeResponse, error) {
	roomTypes, err := hu.hotelRepo.GetRoomTypeByHotelID(ctx, hotelID)
	if err != nil {
		logger.Error(ctx, "Error getting room types by hotel Id", "hotelID", hotelID, "err", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListRoomTypeResponse{}
	resp.RoomTypes = make([]hoteldto.ListRoomType, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		resp.RoomTypes = append(resp.RoomTypes, hoteldto.ListRoomType{
			ID:   roomType.ID,
			Name: roomType.Name,
		})
	}

	return resp, nil
}
