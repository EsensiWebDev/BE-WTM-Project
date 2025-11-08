package booking_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetBookings(ctx context.Context, filter *filter.BookingFilter) ([]entity.Booking, int64, error) {
	db := br.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Booking{})

	// Apply search filter
	if strings.TrimSpace(filter.BookingIDSearch) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.BookingIDSearch)
		query = query.Where("booking_code ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}
	if strings.TrimSpace(filter.GuestNameSearch) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.GuestNameSearch)
		query = query.Joins("JOIN booking_guests bg ON booking_guests.booking_id = bookings.id").
			Where("bg.name ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	// Apply other filters
	if filter.StatusBookingID > 0 {
		query = query.Where("status_booking_id = ?", filter.StatusBookingID)
	}

	if filter.StatusPaymentID > 0 {
		query = query.Where("status_payment_id = ?", filter.StatusPaymentID)
	}

	if filter.AgentID > 0 {
		query = query.Where("agent_id = ?", filter.AgentID)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting facilities", err.Error())
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// Fetch results
	var bookings []model.Booking
	if err := query.
		Preload("StatusBooking").
		Preload("StatusPayment").
		Preload("Agent").
		Find(&bookings).Error; err != nil {
		logger.Error(ctx, "Error fetching bookings", err.Error())
		return nil, 0, err
	}

	// Convert to entity.Booking
	var result []entity.Booking
	if err := utils.CopyStrict(&result, &bookings); err != nil {
		logger.Error(ctx, "Failed to copy bookings to entity", err.Error())
		return nil, 0, err
	}

	return result, total, nil
}
