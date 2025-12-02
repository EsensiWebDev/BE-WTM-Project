package booking_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

var validateColumnSort = map[string]bool{
	"status_booking_id": true,
	"status_payment_id": true,
}

func (br *BookingRepository) GetBookings(ctx context.Context, filter *filter.BookingFilter) ([]entity.Booking, int64, error) {
	db := br.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Booking{})

	// Apply search filter
	if strings.TrimSpace(filter.BookingIDSearch) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.BookingIDSearch)
		query = query.Where("booking_code ILIKE ? ", "%"+safeSearch+"%")
	}

	if strings.TrimSpace(filter.GuestNameSearch) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.GuestNameSearch)
		query = query.Where(`
        EXISTS (
            SELECT 1 FROM booking_details bd
            WHERE bd.booking_id = bookings.id
            AND bd.guest ILIKE ? 
        )`, "%"+safeSearch+"%")
	}

	if filter.BookingStatusID > 0 {
		query = query.Where(`
        EXISTS (
            SELECT 1 FROM booking_details bd
            WHERE bd.booking_id = bookings.id
            AND bd.status_booking_id = ?
        )`, filter.BookingStatusID)
	}

	if filter.PaymentStatusID > 0 {
		query = query.Where(`
        EXISTS (
            SELECT 1 FROM booking_details bd
            WHERE bd.booking_id = bookings.id
            AND bd.status_payment_id = ?
        )`, filter.PaymentStatusID)
	}

	if filter.AgentID > 0 {
		query = query.Where("agent_id = ?", filter.AgentID)
	}

	if filter.ConfirmDateFrom != "" {
		query = query.Where("bookings.confirm_date >= ?", filter.ConfirmDateFrom)
	}

	query = query.Where("bookings.status_booking_id != ?", constant.StatusBookingInCartID)

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

	// Apply sorting
	if filter.Sort != "" {
		if validateColumnSort[filter.Sort] {
			var desc bool
			if strings.TrimSpace(strings.ToLower(filter.Dir)) == "asc" {
				desc = false
			} else {
				desc = true
			}
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: filter.Sort}, Desc: desc})
		}
	}
	query = query.Order("created_at desc")

	// Fetch results
	var bookings []model.Booking
	if err := query.
		Preload("StatusBooking").
		Preload("StatusPayment").
		Preload("Agent").
		Preload("Agent.AgentCompany").
		Preload("Agent.PromoGroup").
		Preload("BookingGuests").
		Preload("BookingDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("BookingDetails.Invoice").
		Preload("BookingDetails.StatusBooking").
		Preload("BookingDetails.StatusPayment").
		Debug().
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
	for i, booking := range bookings {
		for _, guest := range booking.BookingGuests {
			result[i].Guests = append(result[i].Guests, guest.Name)
		}
		result[i].BookingStatus = booking.StatusBooking.Status
		result[i].PaymentStatus = booking.StatusPayment.Status
		for i2, detail := range booking.BookingDetails {
			if detail.Invoice != nil {
				var invoiceEntity entity.DetailInvoice
				if err := json.Unmarshal(detail.Invoice.Detail, &invoiceEntity); err != nil {
					logger.Error(ctx, fmt.Sprintf("Error unmarshalling invoice detail to JSON: %s with detail ID %d", err.Error(), detail.ID), err)
				}
				result[i].BookingDetails[i2].Invoice.DetailInvoice = invoiceEntity
			}
			var detailRoom entity.DetailRoom
			if err := json.Unmarshal(detail.DetailRoom, &detailRoom); err != nil {
				logger.Error(ctx, "Error unmarshalling detail room to JSON", err.Error())
			}
			result[i].BookingDetails[i2].DetailRooms = detailRoom
			result[i].BookingDetails[i2].BookingStatus = detail.StatusBooking.Status
			result[i].BookingDetails[i2].PaymentStatus = detail.StatusPayment.Status
		}
		result[i].AgentName = booking.Agent.FullName
		if booking.Agent.AgentCompany != nil {
			result[i].AgentCompanyName = booking.Agent.AgentCompany.Name
		}
		if booking.Agent.PromoGroup != nil {
			result[i].PromoGroupAgent = booking.Agent.PromoGroup.Name
		}
	}

	return result, total, nil
}
