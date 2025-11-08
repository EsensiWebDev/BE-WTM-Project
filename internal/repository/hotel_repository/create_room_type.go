package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) CreateRoomType(ctx context.Context, roomType *entity.RoomType) (*entity.RoomType, error) {
	db := hr.db.GetTx(ctx)

	var roomTypeModel model.RoomType
	if err := utils.CopyStrict(&roomTypeModel, roomType); err != nil {
		logger.Error(ctx, "Failed to copy room type entity to model", err.Error())
		return nil, err
	}

	roomTypeModel.Hotel = model.Hotel{}
	if err := db.WithContext(ctx).Create(&roomTypeModel).Error; err != nil {
		logger.Error(ctx, "Failed to create room type", err.Error())
		return nil, err
	}

	if err := utils.CopyStrict(roomType, &roomTypeModel); err != nil {
		logger.Error(ctx, "Failed to copy created hotel back to entity", err.Error())
		return nil, err
	}

	return roomType, nil
}
