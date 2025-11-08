package hotel_repository

import (
	"context"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetRoomUnavailableByRoomTypeIDs(ctx context.Context, roomTypeIDs []uint, month time.Time) ([]entity.RoomUnavailable, error) {
	db := hr.db.GetTx(ctx)

	startDate := month
	endDate := startDate.AddDate(0, 1, 0)

	var results []model.RoomUnavailable
	err := db.Model(&model.RoomUnavailable{}).
		Where("room_type_id IN ?", roomTypeIDs).
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("room_type_id ASC, date ASC").
		Find(&results).Error
	if err != nil {
		if hr.db.ErrRecordNotFound(ctx, err) {
			logger.Info(ctx, "No room unavailable data found", "roomTypeIDs", roomTypeIDs, "startDate", startDate)
			return nil, nil
		}
		logger.Error(ctx, "Failed to fetch room unavailable data", err.Error())
		return nil, err
	}

	var entities []entity.RoomUnavailable
	for _, ru := range results {
		var entityRu entity.RoomUnavailable
		if err := utils.CopyPatch(&entityRu, &ru); err != nil {
			logger.Error(ctx, "Failed to copy room unavailable model to entity", err.Error())
			return nil, err
		}
		entities = append(entities, entityRu)
	}

	return entities, nil
}
