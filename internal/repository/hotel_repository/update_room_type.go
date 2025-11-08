package hotel_repository

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
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

	roomPriceWithBreakfast := model.RoomPrice{
		RoomTypeID:  roomType.ID,
		IsBreakfast: true,
		Pax:         roomType.WithBreakfast.Pax,
		Price:       roomType.WithBreakfast.Price,
		IsShow:      roomType.WithBreakfast.IsShow,
	}
	roomPriceWithBreakfast.ID = roomType.WithBreakfast.ID
	roomPriceWithoutBreakfast := model.RoomPrice{
		RoomTypeID:  roomType.ID,
		IsBreakfast: false,
		Price:       roomType.WithoutBreakfast.Price,
		IsShow:      roomType.WithoutBreakfast.IsShow,
	}
	roomPriceWithoutBreakfast.ID = roomType.WithoutBreakfast.ID

	roomPrices := []model.RoomPrice{roomPriceWithBreakfast, roomPriceWithoutBreakfast}

	for _, rp := range roomPrices {
		if rp.ID > 0 {
			// Update existing
			if err := db.WithContext(ctx).
				Model(&model.RoomPrice{}).
				Where("id = ?", rp.ID).
				Updates(map[string]interface{}{
					"price":        rp.Price,
					"is_show":      rp.IsShow,
					"pax":          rp.Pax,
					"is_breakfast": rp.IsBreakfast,
				}).Error; err != nil {
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

	if err := db.WithContext(ctx).
		Debug().
		Unscoped().
		Where("room_type_id = ?", roomType.ID).
		Where("id NOT IN (?)", addIDs).
		Delete(&model.RoomTypeAdditional{}).Error; err != nil {
		logger.Error(ctx, "Failed to delete existing room type additionals", err.Error())
		return err
	}

	return nil
}
