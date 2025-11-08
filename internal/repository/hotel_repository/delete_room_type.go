package hotel_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) DeleteRoomType(ctx context.Context, roomTypeID uint) error {
	db := hr.db.GetTx(ctx)

	if err := db.WithContext(ctx).Delete(&model.RoomType{}, roomTypeID).Error; err != nil {
		logger.Error(ctx, "Failed to delete room type", err.Error())
		return err
	}

	return nil
}
