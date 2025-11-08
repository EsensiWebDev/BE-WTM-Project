package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetBedTypeByRoomTypeID(ctx context.Context, roomTypeID uint) ([]entity.BedType, error) {
	db := hr.db.GetTx(ctx)

	var roomType model.RoomType
	if err := db.WithContext(ctx).
		Preload("BedType").
		First(&roomType, roomTypeID).Error; err != nil {
		logger.Error(ctx, "Error fetching room type with bed types", err.Error())
		return nil, err
	}

	var bedTypesEntities []entity.BedType
	for _, bt := range roomType.BedTypes {
		var btEntity entity.BedType
		if err := utils.CopyPatch(&btEntity, &bt); err != nil {
			logger.Error(ctx, "Failed to copy bed type model to entity", err.Error())
			return nil, err
		}
		bedTypesEntities = append(bedTypesEntities, btEntity)
	}

	return bedTypesEntities, nil
}
