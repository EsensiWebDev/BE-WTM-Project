package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetRoomPriceByID(ctx context.Context, id uint) (*entity.RoomPrice, error) {
	db := hr.db.GetTx(ctx)

	var rp model.RoomPrice
	if err := db.WithContext(ctx).
		Where("id = ?", id).
		Preload("RoomType").
		Preload("RoomType.Hotel").
		Preload("RoomType.BedTypes").
		First(&rp).Error; err != nil {
		return nil, err
	}

	var roomPrice entity.RoomPrice
	if err := utils.CopyStrict(&roomPrice, rp); err != nil {
		return nil, err
	}

	// Map bed types from RoomType
	var bedTypeNames []string
	for _, bedType := range rp.RoomType.BedTypes {
		bedTypeNames = append(bedTypeNames, bedType.Name)
	}
	roomPrice.RoomType.BedTypeNames = bedTypeNames

	return &roomPrice, nil
}
