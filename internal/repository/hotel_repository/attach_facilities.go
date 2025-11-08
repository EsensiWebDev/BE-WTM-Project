package hotel_repository

import (
	"context"
	"gorm.io/gorm"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) AttachFacilities(ctx context.Context, hotelID uint, facilityNames []string) error {
	db := hr.db.GetTx(ctx)

	// Find or create then attach
	var facilities []model.Facility
	for _, name := range facilityNames {
		f := model.Facility{Name: name}
		if err := db.WithContext(ctx).Where("name = ?", name).FirstOrCreate(&f).Error; err != nil {
			logger.Error(ctx, "Failed to find or create facility", err.Error())
			return err
		}
		facilities = append(facilities, f)
	}
	// Association attach
	if err := db.WithContext(ctx).Debug().Model(&model.Hotel{Model: gorm.Model{ID: hotelID}}).
		Association("Facilities").Replace(facilities); err != nil {
		logger.Error(ctx, "Failed to attach facilities to hotel", err.Error())
		return err
	}

	return nil
}
