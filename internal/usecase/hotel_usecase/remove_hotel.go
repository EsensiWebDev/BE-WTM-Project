package hotel_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) RemoveHotel(ctx context.Context, hotelID uint) error {
	if err := hu.hotelRepo.DeleteHotel(ctx, hotelID); err != nil {
		logger.Error(ctx, "Error deleting hotel", "hotelID", hotelID, "err", err.Error())
		return err
	}

	return nil
}
