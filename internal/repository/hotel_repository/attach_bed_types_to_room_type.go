package hotel_repository

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) AttachBedTypesToRoomType(ctx context.Context, roomTypeID uint, bedTypeNames []string) error {
	db := hr.db.GetTx(ctx)

	var bedTypes []model.BedType
	for _, name := range bedTypeNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		var bt model.BedType
		name = utils.CapitalizeWords(name)
		if err := db.Where("name = ?", name).
			Attrs(model.BedType{Name: name}).
			FirstOrCreate(&bt).Error; err != nil {
			logger.Error(ctx, "Failed to find or create bed type", err.Error())
			return err
		}
		bedTypes = append(bedTypes, bt)
	}

	if err := db.Model(&model.RoomType{Model: gorm.Model{ID: roomTypeID}}).
		Association("BedTypes").Replace(bedTypes); err != nil {
		logger.Error(ctx, "Failed to attach bed types to room type", err.Error())
		return err
	}

	return nil
}
