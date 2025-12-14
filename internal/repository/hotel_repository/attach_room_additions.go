package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) AttachRoomAdditions(ctx context.Context, roomTypeID uint, additionals []entity.CustomRoomAdditional) error {
	db := hr.db.GetTx(ctx)
	for _, a := range additionals {
		ra := model.RoomAdditional{Name: a.Name}
		if err := db.Where("name = ?", a.Name).FirstOrCreate(&ra).Error; err != nil {
			logger.Error(ctx, "Failed to create room addition", err.Error())
			return err
		}

		link := model.RoomTypeAdditional{
			RoomTypeID:       roomTypeID,
			RoomAdditionalID: ra.ID,
			Category:         a.Category,
			Price:            a.Price,
			Pax:              a.Pax,
			IsRequired:       a.IsRequired,
		}
		if err := db.Create(&link).Error; err != nil {
			logger.Error(ctx, "Failed to attach room addition to room type", err.Error())
			return err
		}
	}
	return nil
}
