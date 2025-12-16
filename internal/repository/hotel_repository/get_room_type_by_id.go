package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetRoomTypeByID(ctx context.Context, roomTypeID uint) (*entity.RoomType, error) {
	db := hr.db.GetTx(ctx)

	var roomType model.RoomType
	err := db.Where("id = ?", roomTypeID).
		Preload("BedTypes").
		Preload("RoomTypeAdditionals").
		Preload("RoomTypeAdditionals.RoomAdditional").
		Preload("RoomTypePreferences").
		Preload("RoomTypePreferences.OtherPreference").
		Preload("RoomPrices").
		First(&roomType).Error
	if err != nil {
		logger.Error(ctx, "Error fetching room type by Id", err.Error())
		return nil, err
	}

	var roomTypeEntity entity.RoomType
	if err := utils.CopyStrict(&roomTypeEntity, &roomType); err != nil {
		logger.Error(ctx, "Failed to copy room type model to entity", err.Error())
		return nil, err
	}
	var bedTypeNames []string
	for _, bedType := range roomType.BedTypes {
		bedTypeNames = append(bedTypeNames, bedType.Name)
	}
	roomTypeEntity.BedTypeNames = bedTypeNames

	for _, price := range roomType.RoomPrices {
		// Convert Prices JSONB to map
		var pricesMap map[string]float64
		if len(price.Prices) > 0 {
			prices, err := currency.JSONToPrices(price.Prices)
			if err != nil {
				logger.Error(ctx, "Failed to convert prices JSONB to map", err.Error())
				// Fallback to Price field if JSONB conversion fails
				if price.Price > 0 {
					pricesMap = map[string]float64{"IDR": price.Price}
				}
			} else {
				pricesMap = prices
			}
		} else if price.Price > 0 {
			// Fallback: use Price field if Prices JSONB is empty
			pricesMap = map[string]float64{"IDR": price.Price}
		}

		if price.IsBreakfast {
			roomTypeEntity.WithBreakfast = entity.CustomBreakfastWithID{
				ID:     price.ID,
				Pax:    price.Pax,
				Price:  price.Price, // Keep for backward compatibility
				Prices: pricesMap,
				IsShow: price.IsShow,
			}
		} else {
			roomTypeEntity.WithoutBreakfast = entity.CustomBreakfastWithID{
				ID:     price.ID,
				Price:  price.Price, // Keep for backward compatibility
				Prices: pricesMap,
				IsShow: price.IsShow,
			}
		}
	}

	for _, additional := range roomType.RoomTypeAdditionals {
		// Convert Prices JSONB to map for additional services
		var pricesMap map[string]float64
		if len(additional.Prices) > 0 {
			prices, err := currency.JSONToPrices(additional.Prices)
			if err != nil {
				logger.Error(ctx, "Failed to convert additional prices JSONB to map", err.Error())
				// Fallback to Price field if JSONB conversion fails
				if additional.Price != nil && *additional.Price > 0 {
					pricesMap = map[string]float64{"IDR": *additional.Price}
				}
			} else {
				pricesMap = prices
			}
		} else if additional.Price != nil && *additional.Price > 0 {
			// Fallback: use Price field if Prices JSONB is empty
			pricesMap = map[string]float64{"IDR": *additional.Price}
		}

		roomTypeEntity.RoomAdditions = append(roomTypeEntity.RoomAdditions, entity.CustomRoomAdditionalWithID{
			ID:         additional.ID,
			Name:       additional.RoomAdditional.Name,
			Category:   additional.Category,
			Price:      additional.Price, // Keep for backward compatibility
			Prices:     pricesMap,
			Pax:        additional.Pax,
			IsRequired: additional.IsRequired,
		})
	}

	for _, pref := range roomType.RoomTypePreferences {
		roomTypeEntity.OtherPreferences = append(roomTypeEntity.OtherPreferences, entity.CustomOtherPreferenceWithID{
			ID:   pref.ID,
			Name: pref.OtherPreference.Name,
		})
	}

	return &roomTypeEntity, nil
}
