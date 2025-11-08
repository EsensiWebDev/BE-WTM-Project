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

	var bookingDetailAdditional model.BookingDetailAdditional
	if err := utils.CopyStrict(&bookingDetailAdditional, add); err != nil {
		logger.Error(ctx, "Failed to copy booking detail additional entity to model", err.Error())
		return err
	}

	for _, id := range add.BookingDetailIDs {
		bookingDetailAdditional.BookingDetailID = id
		bookingDetailAdditional.ID = 0
		if err := db.WithContext(ctx).Create(&bookingDetailAdditional).Error; err != nil {
			logger.Error(ctx, "Failed to create booking detail additional", err.Error())
			return err
		}
	}

	return nil
}
