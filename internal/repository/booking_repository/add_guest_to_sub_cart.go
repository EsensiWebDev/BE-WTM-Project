package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) AddGuestToSubCart(ctx context.Context, agentID uint, bookingDetailID uint, guest string) error {
	db := br.db.GetTx(ctx)

	// Validasi detail booking dengan detail booking id dan agent id dari booking dan status boooking nya 1
	var exists bool
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Select("count(*) > 0").
		Joins("JOIN bookings ON bookings.id = booking_details.booking_id").
		Where("booking_details.id = ?", bookingDetailID).
		Where("bookings.agent_id = ?", agentID).
		Where("bookings.status_booking_id = ?", constant.StatusBookingInCartID).
		Find(&exists).Error; err != nil {
		logger.Error(ctx, "failed to validate booking detail id with agent id", err.Error())
		return err
	}

	if !exists {
		logger.Error(ctx, "booking detail id not found for the agent in cart status")
		return fmt.Errorf("booking detail id not found for the agent in cart status")
	}

	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("id = ?", bookingDetailID).
		Update("guest", guest).Error; err != nil {
		logger.Error(ctx, "failed to add guest to sub cart", err.Error())
		return err
	}

	return nil
}
