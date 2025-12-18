package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
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

	// Convert Prices JSONB to map
	if len(rp.Prices) > 0 {
		prices, err := currency.JSONToPrices(rp.Prices)
		if err != nil {
			logger.Error(ctx, "Failed to convert prices JSONB to map", err.Error())
			// Fallback to Price field if JSONB conversion fails
			if rp.Price > 0 {
				roomPrice.Prices = map[string]float64{"IDR": rp.Price}
			}
		} else {
			roomPrice.Prices = prices
		}
	} else if rp.Price > 0 {
		// Fallback: use Price field if Prices JSONB is empty
		roomPrice.Prices = map[string]float64{"IDR": rp.Price}
	}

	// Map bed types from RoomType
	var bedTypeNames []string
	for _, bedType := range rp.RoomType.BedTypes {
		bedTypeNames = append(bedTypeNames, bedType.Name)
	}
	roomPrice.RoomType.BedTypeNames = bedTypeNames

	return &roomPrice, nil
}
