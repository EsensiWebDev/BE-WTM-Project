package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
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
		if price.IsBreakfast {
			roomTypeEntity.WithBreakfast = entity.CustomBreakfastWithID{
				ID:     price.ID,
				Pax:    price.Pax,
				Price:  price.Price,
				IsShow: price.IsShow,
			}
		} else {
			roomTypeEntity.WithoutBreakfast = entity.CustomBreakfastWithID{
				ID:     price.ID,
				Price:  price.Price,
				IsShow: price.IsShow,
			}
		}
	}

	for _, additional := range roomType.RoomTypeAdditionals {
		roomTypeEntity.RoomAdditions = append(roomTypeEntity.RoomAdditions, entity.CustomRoomAdditionalWithID{
			ID:         additional.ID,
			Name:       additional.RoomAdditional.Name,
			Category:   additional.Category,
			Price:      additional.Price,
			Pax:        additional.Pax,
			IsRequired: additional.IsRequired,
		})
	}

	return &roomTypeEntity, nil
}
