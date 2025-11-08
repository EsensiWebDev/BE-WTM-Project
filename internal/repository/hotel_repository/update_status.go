package hotel_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) UpdateStatus(ctx context.Context, hotelID uint, statusID uint) error {
	db := hr.db.GetTx(ctx)

	var hotel model.Hotel
	if err := db.WithContext(ctx).First(&hotel, hotelID).Error; err != nil {
		logger.Error(ctx, "Error fetching hotel by Id", err.Error())
		return err
	}

	if hotel.StatusID == statusID {
		logger.Warn(ctx, " Hotel status is already the same as the requested status")
		return nil
	}

	if err := db.WithContext(ctx).
		Model(&model.Hotel{}).
		Where("id = ?", hotelID).
		Update("status_id", statusID).Error; err != nil {
		return err
	}

	return nil
}
