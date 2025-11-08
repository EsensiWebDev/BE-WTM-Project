package hotel_repository

import (
	"context"
	"time"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) DeleteRoomUnavailable(ctx context.Context, roomTypeID uint, month time.Time) error {
	db := hr.db.GetTx(ctx)

	startDate := month
	endDate := startDate.AddDate(0, 1, 0)

	if err := db.WithContext(ctx).Where("date BETWEEN ? AND ?", startDate, endDate).
		Where("room_type_id = ?", roomTypeID).
		Unscoped().Delete(&model.RoomUnavailable{}).Error; err != nil {
		logger.Error(ctx, "Failed to delete room unavailable data", err.Error())
		return err
	}

	return nil
}
