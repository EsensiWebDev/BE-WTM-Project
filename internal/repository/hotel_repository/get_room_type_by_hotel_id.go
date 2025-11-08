package hotel_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetRoomTypeByHotelID(ctx context.Context, hotelID uint) ([]entity.RoomType, error) {
	db := hr.db.GetTx(ctx)
	// Select default fields
	selectFields := []string{"id", "name"}

	// Initialize query with default fields
	query := db.WithContext(ctx).Model(&model.RoomType{}).Select(selectFields)

	// Apply filters
	if hotelID > 0 {
		query = query.Where("hotel_id = ?", hotelID)
	}

	// Execute
	var roomTypes []model.RoomType
	if err := query.Find(&roomTypes).Error; err != nil {
		logger.Error(ctx, "Error fetching room types", err.Error())
		return nil, err
	}

	// Mapping
	var roomTypesEntities []entity.RoomType
	for _, roomType := range roomTypes {
		var roomTypeEntity entity.RoomType
		if err := utils.CopyPatch(&roomTypeEntity, &roomType); err != nil {
			logger.Error(ctx, "Failed to copy room type model to entity", err.Error())
			return nil, err
		}
		roomTypesEntities = append(roomTypesEntities, roomTypeEntity)
	}

	return roomTypesEntities, nil
}
