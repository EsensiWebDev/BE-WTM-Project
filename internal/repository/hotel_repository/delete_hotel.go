package hotel_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) DeleteHotel(ctx context.Context, hotelID uint) error {
	db := hr.db.GetTx(ctx)

	if err := db.WithContext(ctx).Delete(&model.Hotel{}, hotelID).Error; err != nil {
		logger.Error(ctx, "Failed to delete hotel", err.Error())
		return err
	}

	return nil
}
