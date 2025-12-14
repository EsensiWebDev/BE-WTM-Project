package booking_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetBookingByID(ctx context.Context, bookingID uint) (*entity.Booking, error) {
	db := br.db.GetTx(ctx)

	var booking model.Booking
	if err := db.
		Preload("BookingDetails").
		Preload("BookingGuests").
		Where("id = ?", bookingID).
		First(&booking).Error; err != nil {
		logger.Error(ctx, "failed to get booking by Id", err.Error())
		return nil, err
	}

	var result entity.Booking
	if err := utils.CopyStrict(&result, &booking); err != nil {
		logger.Error(ctx, "failed to copy booking model to entity", err.Error())
		return nil, err
	}

	// Map guests from model to entity (entity only has []string, but we need full details for email)
	// For now, just map names - full guest details will be accessed from model in usecase
	var guests []string
	for _, guest := range booking.BookingGuests {
		guests = append(guests, guest.Name)
	}
	result.Guests = guests

	return &result, nil

}
