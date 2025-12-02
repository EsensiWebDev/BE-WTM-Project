package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) GetSubBookingIDs(ctx context.Context, agentID uint, bookingCode string) ([]string, error) {
	db := br.db.GetTx(ctx)

	var subBookingIDs []string

	err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Select("booking_details.sub_booking_id").
		Joins("JOIN bookings ON bookings.id = booking_details.booking_id").
		Where("bookings.agent_id = ? AND bookings.booking_code = ?", agentID, bookingCode).
		Order("booking_details.created_at DESC").
		Pluck("booking_details.sub_booking_id", &subBookingIDs).Error

	if err != nil {
		logger.Error(ctx, "failed to fetch sub booking ids", err.Error())
		return nil, err
	}

	return subBookingIDs, nil
}
