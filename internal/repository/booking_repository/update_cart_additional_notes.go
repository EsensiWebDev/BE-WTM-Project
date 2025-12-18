package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

// UpdateCartAdditionalNotes updates booking_details.additional_notes for a single
// booking_detail (cart item), ensuring it belongs to the specified agent.
func (br *BookingRepository) UpdateCartAdditionalNotes(ctx context.Context, agentID uint, bookingDetailID uint, additionalNotes string) error {
	db := br.db.GetTx(ctx)

	// Ensure the booking_detail belongs to the agent's booking
	var booking model.Booking
	if err := db.WithContext(ctx).
		Model(&model.Booking{}).
		Joins("JOIN booking_details bd ON bd.booking_id = bookings.id").
		Where("bd.id = ? AND bookings.agent_id = ?", bookingDetailID, agentID).
		Select("bookings.id").
		First(&booking).Error; err != nil {
		logger.Error(ctx, "booking detail not found or not owned by agent for updating additional notes", err.Error())
		return fmt.Errorf("unauthorized or not found")
	}

	// Perform the update on booking_details
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("id = ?", bookingDetailID).
		Update("additional_notes", additionalNotes).Error; err != nil {
		logger.Error(ctx, "failed to update additional notes on booking detail", err.Error())
		return fmt.Errorf("failed to update additional notes")
	}

	return nil
}



