package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) GetOrCreateCartID(ctx context.Context, agentID uint) (uint, error) {
	db := br.db.GetTx(ctx)

	var booking model.Booking
	if err := db.WithContext(ctx).
		Where("agent_id = ? AND status_booking_id = ?", agentID, 1).
		First(&booking).Error; err != nil {
		if br.db.ErrRecordNotFound(ctx, err) {
			bookingCode, err := br.GenerateCode(ctx, "booking_codes", "BK")
			if err != nil {
				logger.Error(ctx, "failed to generate booking code", "error", err)
				return 0, err
			}

			if bookingCode == "" {
				logger.Error(ctx, "failed to generate booking code", "error", err)
				return 0, err
			}

			booking = model.Booking{
				AgentID:         agentID,
				StatusBookingID: constant.StatusBookingInCartID,
				StatusPaymentID: constant.StatusPaymentUnpaidID,
				BookingCode:     bookingCode,
			}
			if err := db.WithContext(ctx).Create(&booking).Error; err != nil {
				logger.Error(ctx, "Error creating booking for cart", "error", err)
				return 0, err
			}
		} else {
			logger.Error(ctx, "Error fetching booking for cart", "error", err)
			return 0, err
		}
	}
	return booking.ID, nil
}
