package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// DeleteAllGuestsFromBooking deletes all guests from a booking
// This is typically called after successful checkout to clean up contact details
func (br *BookingRepository) DeleteAllGuestsFromBooking(ctx context.Context, bookingID uint) error {
	db := br.db.GetTx(ctx)

	if err := db.WithContext(ctx).
		Unscoped().
		Where("booking_id = ?", bookingID).
		Delete(&model.BookingGuest{}).Error; err != nil {
		logger.Error(ctx, "failed to delete all guests from booking", err.Error())
		return err
	}

	return nil
}
