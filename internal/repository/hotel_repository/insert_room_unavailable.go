package hotel_repository

import (
	"context"
	"time"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) InsertRoomUnavailable(ctx context.Context, roomTypeID uint, unavailableDates []time.Time) error {
	db := hr.db.GetTx(ctx)

	if len(unavailableDates) == 0 {
		logger.Info(ctx, "No unavailable dates provided for room type", "roomTypeID", roomTypeID)
		return nil
	}

	var roomUnavailable []model.RoomUnavailable
	for _, date := range unavailableDates {
		roomUnavailable = append(roomUnavailable, model.RoomUnavailable{
			RoomTypeID: roomTypeID,
			Date:       &date,
		})
	}

	if err := db.WithContext(ctx).Create(&roomUnavailable).Error; err != nil {
		logger.Error(ctx, "Failed to insert room unavailable data", err.Error())
		return err
	}

	return nil
}
