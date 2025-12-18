package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetRoomTypeAdditionalsByIDs(ctx context.Context, ids []uint) ([]entity.RoomTypeAdditional, error) {
	db := hr.db.GetTx(ctx)

	if len(ids) == 0 {
		return nil, nil
	}
	var additionals []model.RoomTypeAdditional
	if err := db.WithContext(ctx).
		Preload("RoomAdditional").
		Where("id IN ?", ids).
		Find(&additionals).Error; err != nil {
		if hr.db.ErrRecordNotFound(ctx, err) {
			logger.Error(ctx, "Not found")
			return nil, nil
		}
		return nil, err
	}

	// Convert model.RoomTypeAdditional to entity.RoomTypeAdditional
	var additionalsEntity []entity.RoomTypeAdditional
	if err := utils.CopyStrict(&additionalsEntity, &additionals); err != nil {
		return nil, err
	}

	// Convert Prices JSONB to map for each additional
	for i, add := range additionals {
		var pricesMap map[string]float64
		if len(add.Prices) > 0 {
			prices, err := currency.JSONToPrices(add.Prices)
			if err != nil {
				logger.Error(ctx, "Failed to convert additional prices JSONB to map", err.Error())
				// Fallback to Price field if JSONB conversion fails
				if add.Price != nil && *add.Price > 0 {
					pricesMap = map[string]float64{"IDR": *add.Price}
				} else {
					pricesMap = make(map[string]float64)
				}
			} else {
				pricesMap = prices
			}
		} else if add.Price != nil && *add.Price > 0 {
			// Fallback: use Price field if Prices JSONB is empty
			pricesMap = map[string]float64{"IDR": *add.Price}
		} else {
			// Initialize empty map to avoid nil
			pricesMap = make(map[string]float64)
		}
		additionalsEntity[i].Prices = pricesMap
	}

	return additionalsEntity, nil
}
