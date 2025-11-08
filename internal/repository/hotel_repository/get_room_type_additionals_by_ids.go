package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
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

	return additionalsEntity, nil
}
