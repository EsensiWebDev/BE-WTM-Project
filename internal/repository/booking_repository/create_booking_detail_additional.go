package booking_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) CreateBookingDetailAdditional(ctx context.Context, add *entity.BookingDetailAdditional) error {
	db := br.db.GetTx(ctx)

	for _, id := range add.BookingDetailIDs {
		// Create a fresh model instance for each BookingDetailID so that
		// the embedded ExternalID starts empty and the BeforeCreate hook
		// can safely generate a new unique value per row.
		var bookingDetailAdditional model.BookingDetailAdditional
		if err := utils.CopyStrict(&bookingDetailAdditional, add); err != nil {
			logger.Error(ctx, "Failed to copy booking detail additional entity to model", err.Error())
			return err
		}

		bookingDetailAdditional.BookingDetailID = id
		bookingDetailAdditional.ID = 0
		if err := db.WithContext(ctx).Create(&bookingDetailAdditional).Error; err != nil {
			logger.Error(ctx, "Failed to create booking detail additional", err.Error())
			return err
		}
	}

	return nil
}
