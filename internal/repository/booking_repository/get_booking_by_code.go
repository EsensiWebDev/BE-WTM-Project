package booking_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetBookingByCode(ctx context.Context, code string) (*entity.Booking, error) {
	db := br.db.GetTx(ctx)

	var booking model.Booking
	if err := db.
		Preload("Agent").
		Preload("Agent.AgentCompany").
		Preload("BookingGuests").
		Preload("BookingDetails").
		Preload("BookingDetails.RoomType").
		Preload("BookingDetails.RoomType.Hotel").
		Where("booking_code = ?", code).
		First(&booking).Error; err != nil {
		logger.Error(ctx, "failed to get booking by Id", err.Error())
		return nil, err
	}

	var result entity.Booking
	if err := utils.CopyStrict(&result, &booking); err != nil {
		logger.Error(ctx, "failed to copy booking model to entity", err.Error())
		return nil, err
	}
	result.AgentCompanyName = "-"
	if booking.Agent.AgentCompany != nil {
		result.AgentCompanyName = booking.Agent.AgentCompany.Name
	}
	result.AgentPhoneNumber = "-"
	if booking.Agent.Phone != "" {
		result.AgentPhoneNumber = booking.Agent.Phone
	}
	if len(booking.BookingGuests) > 0 {
		result.Guests = make([]string, len(booking.BookingGuests))
		for i, guest := range booking.BookingGuests {
			result.Guests[i] = guest.Name
		}
	}

	return &result, nil

}
