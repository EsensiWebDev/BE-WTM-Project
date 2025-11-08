package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateBookingStatus(ctx context.Context, bookingID uint, statusBookingID uint) error {
	db := br.db.GetTx(ctx)

	err := db.Model(&model.Booking{}).
		Where("id = ?", bookingID).
		Update("status_booking_id", statusBookingID).Error
	if err != nil {
		logger.Error(ctx, "failed to update booking status: ", err.Error())
		return fmt.Errorf("failed to update booking status: %s", err.Error())
	}

	return nil
}
