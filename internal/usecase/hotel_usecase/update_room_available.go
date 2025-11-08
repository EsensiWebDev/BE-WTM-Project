package hotel_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UpdateRoomAvailable(ctx context.Context, req *hoteldto.UpdateRoomAvailableRequest) error {
	return hu.dbTransaction.WithTransaction(ctx, func(txCtx context.Context) error {
		monthTime, err := time.Parse("2006-01", req.Month)
		if err != nil {
			logger.Error(ctx, "Failed to parse month", err.Error())
			return err
		}

		for _, data := range req.Data {
			unavailableDates, err := buildUnavailableDates(monthTime, data.RoomAvailable)
			if err != nil {
				logger.Error(ctx, "Failed to build unavailable dates", err.Error())
				return fmt.Errorf("room_type_id %d: build error: %s", data.RoomTypeID, err.Error())
			}

			if err := hu.hotelRepo.DeleteRoomUnavailable(txCtx, data.RoomTypeID, monthTime); err != nil {
				logger.Error(ctx, "Failed to delete room unavailability", "roomTypeID", data.RoomTypeID, err.Error())
				return fmt.Errorf("room_type_id %d: delete error: %s", data.RoomTypeID, err.Error())
			}

			if len(unavailableDates) > 0 {
				if err := hu.hotelRepo.InsertRoomUnavailable(txCtx, data.RoomTypeID, unavailableDates); err != nil {
					logger.Error(ctx, "Failed to insert room unavailability", "roomTypeID", data.RoomTypeID, err.Error())
					return fmt.Errorf("room_type_id %d: insert error: %s", data.RoomTypeID, err.Error())
				}
			}
		}

		return nil
	})
}

func buildUnavailableDates(month time.Time, days []hoteldto.DataAvailable) ([]time.Time, error) {
	var unavailable []time.Time
	for _, day := range days {
		if !day.Available {
			date := time.Date(month.Year(), month.Month(), day.Day, 0, 0, 0, 0, time.UTC)
			unavailable = append(unavailable, date)
		}
	}
	return unavailable, nil
}
