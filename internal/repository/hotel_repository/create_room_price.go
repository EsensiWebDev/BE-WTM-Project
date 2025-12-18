package hotel_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) CreateRoomPrice(ctx context.Context, roomTypeID uint, dto *entity.CustomBreakfast, isBreakfast bool) error {
	db := hr.db.GetTx(ctx)

	// Convert Prices map to JSONB
	var pricesJSON []byte
	var err error
	if dto.Prices != nil && len(dto.Prices) > 0 {
		// Validate prices
		if err := currency.ValidatePrices(dto.Prices); err != nil {
			logger.Error(ctx, "Invalid prices", err.Error())
			return err
		}
		pricesJSON, err = currency.PricesToJSON(dto.Prices)
		if err != nil {
			logger.Error(ctx, "Failed to convert prices to JSON", err.Error())
			return err
		}
	} else if dto.Price > 0 {
		// Fallback: convert single Price to Prices JSONB with IDR
		pricesJSON, err = currency.PricesToJSON(map[string]float64{"IDR": dto.Price})
		if err != nil {
			logger.Error(ctx, "Failed to convert price to JSON", err.Error())
			return err
		}
	} else {
		logger.Error(ctx, "Price or Prices must be provided")
		return errors.New("price or prices must be provided")
	}

	rp := &model.RoomPrice{
		RoomTypeID:  roomTypeID,
		IsBreakfast: isBreakfast,
		Pax:         dto.Pax,
		Price:       dto.Price, // Keep for backward compatibility
		Prices:      pricesJSON,
		IsShow:      dto.IsShow,
	}

	if err := db.Create(rp).Error; err != nil {
		logger.Error(ctx, "Failed to create room price", err.Error())
		return err
	}

	return nil
}
