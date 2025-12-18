package hotel_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) UpdateRoomType(ctx context.Context, roomType *entity.RoomType) error {
	db := hr.db.GetTx(ctx)

	if roomType.ID == 0 {
		return errors.New("invalid room type ID")
	}

	var roomTypeModel model.RoomType
	if err := utils.CopyStrict(&roomTypeModel, roomType); err != nil {
		logger.Error(ctx, "Failed to copy room type entity to model", err.Error())
		return err
	}

	roomTypeModel.BedTypes = nil
	roomTypeModel.PromoRoomTypes = nil
	roomTypeModel.RoomTypeAdditionals = nil
	roomTypeModel.RoomPrices = nil

	if err := db.WithContext(ctx).Model(&model.RoomType{}).
		Where("id = ?", roomType.ID).
		Updates(&roomTypeModel).Error; err != nil {
		logger.Error(ctx, "failed to update room type", err.Error())
		return err
	}

	// Helper function to convert Prices map to JSONB
	convertPricesToJSONB := func(prices map[string]float64, fallbackPrice float64) ([]byte, error) {
		if prices != nil && len(prices) > 0 {
			// Validate prices
			if err := currency.ValidatePrices(prices); err != nil {
				return nil, err
			}
			return currency.PricesToJSON(prices)
		} else if fallbackPrice > 0 {
			// Fallback: convert single Price to Prices JSONB with IDR
			return currency.PricesToJSON(map[string]float64{"IDR": fallbackPrice})
		}
		return nil, errors.New("price or prices must be provided")
	}

	// Prepare room price with breakfast
	withBreakfastPricesJSON, err := convertPricesToJSONB(roomType.WithBreakfast.Prices, roomType.WithBreakfast.Price)
	if err != nil {
		logger.Error(ctx, "Failed to convert with breakfast prices to JSONB", err.Error())
		return err
	}
	roomPriceWithBreakfast := model.RoomPrice{
		RoomTypeID:  roomType.ID,
		IsBreakfast: true,
		Pax:         roomType.WithBreakfast.Pax,
		Price:       roomType.WithBreakfast.Price, // Keep for backward compatibility
		Prices:      withBreakfastPricesJSON,
		IsShow:      roomType.WithBreakfast.IsShow,
	}
	roomPriceWithBreakfast.ID = roomType.WithBreakfast.ID

	// Prepare room price without breakfast
	withoutBreakfastPricesJSON, err := convertPricesToJSONB(roomType.WithoutBreakfast.Prices, roomType.WithoutBreakfast.Price)
	if err != nil {
		logger.Error(ctx, "Failed to convert without breakfast prices to JSONB", err.Error())
		return err
	}
	roomPriceWithoutBreakfast := model.RoomPrice{
		RoomTypeID:  roomType.ID,
		IsBreakfast: false,
		Price:       roomType.WithoutBreakfast.Price, // Keep for backward compatibility
		Prices:      withoutBreakfastPricesJSON,
		IsShow:      roomType.WithoutBreakfast.IsShow,
	}
	roomPriceWithoutBreakfast.ID = roomType.WithoutBreakfast.ID

	roomPrices := []model.RoomPrice{roomPriceWithBreakfast, roomPriceWithoutBreakfast}

	for _, rp := range roomPrices {
		if rp.ID > 0 {
			// Update existing
			updateMap := map[string]interface{}{
				"price":        rp.Price,
				"prices":       rp.Prices,
				"is_show":      rp.IsShow,
				"pax":          rp.Pax,
				"is_breakfast": rp.IsBreakfast,
			}
			if err := db.WithContext(ctx).
				Model(&model.RoomPrice{}).
				Where("id = ?", rp.ID).
				Updates(updateMap).Error; err != nil {
				return err
			}
		} else {
			// Insert new
			if err := db.WithContext(ctx).Create(&rp).Error; err != nil {
				return err
			}
		}
	}

	var addIDs []uint
	for _, add := range roomType.RoomAdditions {
		addIDs = append(addIDs, add.ID)
	}

	// Delete RoomTypeAdditional records that are no longer attached to this room type,
	// regardless of whether they have been used in historical bookings.
	deleteQuery := db.WithContext(ctx).
		Unscoped().
		Where("room_type_id = ?", roomType.ID)

	if len(addIDs) > 0 {
		deleteQuery = deleteQuery.Where("id NOT IN (?)", addIDs)
	}

	if err := deleteQuery.Delete(&model.RoomTypeAdditional{}).Error; err != nil {
		logger.Error(ctx, "Failed to delete existing room type additionals", err.Error())
		return err
	}

	// Clean up orphan RoomAdditional records that are no longer referenced
	var orphanAdditionalIDs []uint
	if err := db.WithContext(ctx).
		Model(&model.RoomAdditional{}).
		Joins("LEFT JOIN room_type_additionals rta ON rta.room_additional_id = room_additionals.id").
		Where("rta.id IS NULL").
		Pluck("room_additionals.id", &orphanAdditionalIDs).Error; err != nil {
		logger.Error(ctx, "Failed to find orphan room additionals", err.Error())
		return err
	}

	if len(orphanAdditionalIDs) > 0 {
		if err := db.WithContext(ctx).
			Unscoped().
			Where("id IN (?)", orphanAdditionalIDs).
			Delete(&model.RoomAdditional{}).Error; err != nil {
			logger.Error(ctx, "Failed to delete orphan room additionals", err.Error())
			return err
		}
	}

	return nil
}
