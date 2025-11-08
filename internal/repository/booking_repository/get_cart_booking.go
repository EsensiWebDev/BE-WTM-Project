package booking_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetCartBooking(ctx context.Context, agentID uint) (*entity.Booking, error) {
	db := br.db.GetTx(ctx)

	var booking model.Booking
	if err := db.WithContext(ctx).
		Where("agent_id = ? AND status_booking_id = ?", agentID, constant.StatusBookingInCartID).
		Preload("BookingDetails").
		Preload("BookingDetails.BookingDetailsAdditional").
		Preload("BookingGuests").
		First(&booking).Error; err != nil {
		if br.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "No cart booking found for agent Id", agentID)
			return nil, nil // No cart booking found
		}
		logger.Error(ctx, "Error finding cart booking for agent Id", agentID, err.Error())
		return nil, err // Other error
	}

	var bookingEntity entity.Booking
	if err := utils.CopyStrict(&bookingEntity, &booking); err != nil {
		logger.Error(ctx, "Error copying booking model to entity", err.Error())
		return nil, err // Error copying model to entity
	}

	var guests []string
	for _, guest := range booking.BookingGuests {
		guests = append(guests, guest.Name)
	}
	bookingEntity.Guests = guests

	for i, detail := range booking.BookingDetails {
		var detailPromo entity.DetailPromo
		if err := json.Unmarshal(detail.DetailPromo, &detailPromo); err != nil {
			logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
		}
		bookingEntity.BookingDetails[i].DetailPromos = detailPromo

		var detailRoom entity.DetailRoom
		if err := json.Unmarshal(detail.DetailRoom, &detailRoom); err != nil {
			logger.Error(ctx, "Error marshalling room detail to JSON", err.Error())
		}
		bookingEntity.BookingDetails[i].DetailRooms = detailRoom
	}

	return &bookingEntity, nil
}
