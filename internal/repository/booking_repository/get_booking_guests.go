package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// GetBookingGuests retrieves full guest details for a booking
func (br *BookingRepository) GetBookingGuests(ctx context.Context, bookingID uint) ([]model.BookingGuest, error) {
	db := br.db.GetTx(ctx)

	var guests []model.BookingGuest
	if err := db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Find(&guests).Error; err != nil {
		logger.Error(ctx, "failed to get booking guests", err.Error())
		return nil, err
	}

	return guests, nil
}
