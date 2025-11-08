package hotel_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) RemoveRoomType(ctx context.Context, roomTypeID uint) error {
	if err := hu.hotelRepo.DeleteRoomType(ctx, roomTypeID); err != nil {
		logger.Error(ctx, "Error deleting room type by Id", "roomTypeID", roomTypeID, "err", err.Error())
		return err
	}

	return nil
}
