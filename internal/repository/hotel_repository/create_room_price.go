package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) CreateRoomPrice(ctx context.Context, roomTypeID uint, dto *entity.CustomBreakfast, isBreakfast bool) error {
	db := hr.db.GetTx(ctx)

	rp := &model.RoomPrice{
		RoomTypeID:  roomTypeID,
		IsBreakfast: isBreakfast,
		Pax:         dto.Pax,
		Price:       dto.Price,
		IsShow:      dto.IsShow,
	}

	if err := db.Create(rp).Error; err != nil {
		logger.Error(ctx, "Failed to create room price", err.Error())
		return err
	}

	return nil
}
