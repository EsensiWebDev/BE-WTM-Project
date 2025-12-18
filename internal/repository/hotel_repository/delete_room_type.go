package hotel_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) DeleteRoomType(ctx context.Context, roomTypeID uint) error {
	db := hr.db.GetTx(ctx)

	// First, delete RoomTypeAdditional records for this room type
	if err := db.WithContext(ctx).
		Unscoped().
		Where("room_type_id = ?", roomTypeID).
		Delete(&model.RoomTypeAdditional{}).Error; err != nil {
		logger.Error(ctx, "Failed to delete room type additionals before deleting room type", err.Error())
		return err
	}

	// Clean up orphan RoomAdditional records that are no longer referenced
	var orphanAdditionalIDs []uint
	if err := db.WithContext(ctx).
		Model(&model.RoomAdditional{}).
		Joins("LEFT JOIN room_type_additionals rta ON rta.room_additional_id = room_additionals.id").
		Where("rta.id IS NULL").
		Pluck("room_additionals.id", &orphanAdditionalIDs).Error; err != nil {
		logger.Error(ctx, "Failed to find orphan room additionals on delete room type", err.Error())
		return err
	}

	if len(orphanAdditionalIDs) > 0 {
		if err := db.WithContext(ctx).
			Unscoped().
			Where("id IN (?)", orphanAdditionalIDs).
			Delete(&model.RoomAdditional{}).Error; err != nil {
			logger.Error(ctx, "Failed to delete orphan room additionals on delete room type", err.Error())
			return err
		}
	}

	// Finally, delete the RoomType itself
	if err := db.WithContext(ctx).Delete(&model.RoomType{}, roomTypeID).Error; err != nil {
		logger.Error(ctx, "Failed to delete room type", err.Error())
		return err
	}

	return nil
}
