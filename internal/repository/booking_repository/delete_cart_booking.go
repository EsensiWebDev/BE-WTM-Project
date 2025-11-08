package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
)

func (br *BookingRepository) DeleteCartBooking(ctx context.Context, agentID uint, bookingDetailID uint) error {
	db := br.db.GetTx(ctx)

	// Delete the booking detail by Id
	if err := db.WithContext(ctx).
		Where("agent_id = ? AND id = ?", agentID, bookingDetailID).
		Delete(&model.BookingDetail{}).Error; err != nil {
		return err
	}

	// Check if there are any remaining booking details for the agent
	var count int64
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("agent_id = ?", agentID).
		Count(&count).Error; err != nil {
		return err
	}

	// If no booking details left, delete the cart booking
	if count == 0 {
		if err := db.WithContext(ctx).
			Where("agent_id = ?", agentID).
			Delete(&model.Booking{}).Error; err != nil {
			return err
		}
	}

	return nil
}
