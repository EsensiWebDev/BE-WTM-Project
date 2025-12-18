package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// UpdateAdminNotes updates booking_details.admin_notes for a single
// booking_detail identified by sub_booking_id.
func (br *BookingRepository) UpdateAdminNotes(ctx context.Context, subBookingID string, adminNotes string) error {
	db := br.db.GetTx(ctx)

	// Check if booking detail exists
	var bookingDetail model.BookingDetail
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("sub_booking_id = ?", subBookingID).
		First(&bookingDetail).Error; err != nil {
		logger.Error(ctx, "booking detail not found for updating admin notes", err.Error())
		return fmt.Errorf("booking detail not found")
	}

	// Perform the update on booking_details
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("sub_booking_id = ?", subBookingID).
		Update("admin_notes", adminNotes).Error; err != nil {
		logger.Error(ctx, "failed to update admin notes on booking detail", err.Error())
		return fmt.Errorf("failed to update admin notes")
	}

	return nil
}
