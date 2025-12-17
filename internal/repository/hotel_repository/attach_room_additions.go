package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
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

		// Convert Prices map to JSONB
		var pricesJSON []byte
		if len(a.Prices) > 0 {
			// Validate prices
			if err := currency.ValidatePrices(a.Prices); err != nil {
				logger.Error(ctx, "Invalid prices for room addition", err.Error())
				return err
			}
			var err error
			pricesJSON, err = currency.PricesToJSON(a.Prices)
			if err != nil {
				logger.Error(ctx, "Failed to convert prices to JSON", err.Error())
				return err
			}
		} else if a.Price != nil && *a.Price > 0 {
			// Fallback: convert single Price to Prices JSONB with IDR
			var err error
			pricesJSON, err = currency.PricesToJSON(map[string]float64{"IDR": *a.Price})
			if err != nil {
				logger.Error(ctx, "Failed to convert price to JSON", err.Error())
				return err
			}
		}

		// Check if RoomTypeAdditional already exists for this room type and addition
		var existingLink model.RoomTypeAdditional
		err := db.Where("room_type_id = ? AND room_additional_id = ?", roomTypeID, ra.ID).First(&existingLink).Error

		if err == nil {
			// Update existing record - merge prices with existing ones
			var mergedPrices map[string]float64

			// Get existing prices from database
			if len(existingLink.Prices) > 0 {
				existingPrices, err := currency.JSONToPrices(existingLink.Prices)
				if err != nil {
					logger.Error(ctx, "Failed to convert existing prices JSONB to map", err.Error())
					// Fallback: use Price field if JSONB conversion fails
					if existingLink.Price != nil && *existingLink.Price > 0 {
						mergedPrices = map[string]float64{"IDR": *existingLink.Price}
					} else {
						mergedPrices = make(map[string]float64)
					}
				} else {
					mergedPrices = existingPrices
				}
			} else if existingLink.Price != nil && *existingLink.Price > 0 {
				// Fallback: use Price field if Prices JSONB is empty
				mergedPrices = map[string]float64{"IDR": *existingLink.Price}
			} else {
				mergedPrices = make(map[string]float64)
			}

			// Merge new prices with existing prices (new currencies added, existing currencies updated)
			if len(a.Prices) > 0 {
				for currency, price := range a.Prices {
					mergedPrices[currency] = price
				}
			} else if a.Price != nil && *a.Price > 0 {
				// Fallback: if only Price is provided, update/add IDR
				mergedPrices["IDR"] = *a.Price
			}

			// Validate merged prices
			if len(mergedPrices) > 0 {
				if err := currency.ValidatePrices(mergedPrices); err != nil {
					logger.Error(ctx, "Invalid merged prices for room addition", err.Error())
					return err
				}
			}

			// Convert merged prices to JSONB
			var finalPricesJSON []byte
			if len(mergedPrices) > 0 {
				var err error
				finalPricesJSON, err = currency.PricesToJSON(mergedPrices)
				if err != nil {
					logger.Error(ctx, "Failed to convert merged prices to JSON", err.Error())
					return err
				}
			}

			// Update deprecated Price field with IDR price from merged prices for backward compatibility
			var updatedPrice *float64
			if idrPrice, exists := mergedPrices["IDR"]; exists {
				updatedPrice = &idrPrice
			} else if a.Price != nil {
				updatedPrice = a.Price
			}

			// Update existing record with merged prices
			updateMap := map[string]interface{}{
				"category":    a.Category,
				"price":       updatedPrice,
				"prices":      finalPricesJSON,
				"pax":         a.Pax,
				"is_required": a.IsRequired,
			}
			if err := db.Model(&existingLink).Updates(updateMap).Error; err != nil {
				logger.Error(ctx, "Failed to update room addition to room type", err.Error())
				return err
			}
		} else {
			// Create new record
			link := model.RoomTypeAdditional{
				RoomTypeID:       roomTypeID,
				RoomAdditionalID: ra.ID,
				Category:         a.Category,
				Price:            a.Price,
				Prices:           pricesJSON,
				Pax:              a.Pax,
				IsRequired:       a.IsRequired,
			}
			if err := db.Create(&link).Error; err != nil {
				logger.Error(ctx, "Failed to attach room addition to room type", err.Error())
				return err
			}
		}
	}
	return nil
}
