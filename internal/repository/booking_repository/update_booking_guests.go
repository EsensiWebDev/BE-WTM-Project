package booking_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateBookingGuests(ctx context.Context, bookingID uint, guests []string) error {
	db := br.db.GetTx(ctx)

	// Step 1: Delete existing guests
	if err := db.WithContext(ctx).
		Unscoped().
		Where("booking_id = ?", bookingID).
		Delete(&model.BookingGuest{}).Error; err != nil {
		logger.Error(ctx, "failed to delete old booking guests: ", err)
		return err
	}

	// Step 2: Insert new guests
	var newGuests []model.BookingGuest
	for _, name := range guests {
		if strings.TrimSpace(name) == "" {
			continue
		}
		newGuests = append(newGuests, model.BookingGuest{
			BookingID: bookingID,
			Name:      name,
		})
	}

	if len(newGuests) == 0 {
		logger.Info(ctx, "no new guests to update")
		return nil
	}

	if err := db.WithContext(ctx).Create(&newGuests).Error; err != nil {
		logger.Error(ctx, "failed to insert new booking guests: ", err)
		return err
	}

	return nil
}
