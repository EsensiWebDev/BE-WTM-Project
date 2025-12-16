package booking_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var validateColumnSort = map[string]bool{
	"status_booking_id": true,
	"status_payment_id": true,
}

func (br *BookingRepository) GetBookings(ctx context.Context, filter *filter.BookingFilter) ([]entity.Booking, int64, error) {
	db := br.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Booking{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where(db.Where("booking_code ILIKE ? ", "%"+safeSearch+"%").
			Or(`
        EXISTS (
            SELECT 1 FROM booking_details bd
            WHERE bd.booking_id = bookings.id
            AND bd.guest ILIKE ? 
        )`, "%"+safeSearch+"%").
			Or(`
        EXISTS (
            SELECT 1 FROM users u
            JOIN agent_companies ac ON u.agent_company_id = ac.id
            WHERE u.id = bookings.agent_id
            AND (u.full_name ILIKE ? OR ac.name ILIKE ?) 
        )`, "%"+safeSearch+"%", "%"+safeSearch+"%"))
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

	if filter.ConfirmDateFrom != nil {
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
		Preload("BookingDetails.BookingDetailsAdditional").
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
			// Unmarshal promo snapshot from booking_details.detail_promo (JSONB) into entity.DetailPromo
			if len(detail.DetailPromo) > 0 {
				var promoSnap entity.DetailPromo
				if err := json.Unmarshal(detail.DetailPromo, &promoSnap); err != nil {
					logger.Error(ctx, fmt.Sprintf("Error unmarshalling detail promo to JSON: %s with detail ID %d", err.Error(), detail.ID), err)
				} else {
					result[i].BookingDetails[i2].DetailPromos = promoSnap
				}
			}

			if detail.Invoice != nil {
				var invoiceEntity entity.DetailInvoice
				if err := json.Unmarshal(detail.Invoice.Detail, &invoiceEntity); err != nil {
					logger.Error(ctx, fmt.Sprintf("Error unmarshalling invoice detail to JSON: %s with detail ID %d", err.Error(), detail.ID), err)
				}
				// Backfill promo into invoice detail from booking detail snapshot if missing
				if invoiceEntity.Promo.Name == "" && result[i].BookingDetails[i2].DetailPromos.Name != "" {
					invoiceEntity.Promo = result[i].BookingDetails[i2].DetailPromos
				}
				result[i].BookingDetails[i2].Invoice.DetailInvoice = invoiceEntity
			}

			var detailRoom entity.DetailRoom
			if err := json.Unmarshal(detail.DetailRoom, &detailRoom); err != nil {
				logger.Error(ctx, "Error unmarshalling detail room to JSON", err.Error())
			}
			result[i].BookingDetails[i2].DetailRooms = detailRoom

			// Map additional services snapshot from booking_detail_additionals
			// Copy full BookingDetailsAdditional objects (for detailed info)
			if len(detail.BookingDetailsAdditional) > 0 {
				// Copy the full objects for detailed information
				additionalEntities := make([]entity.BookingDetailAdditional, 0, len(detail.BookingDetailsAdditional))
				names := make([]string, 0, len(detail.BookingDetailsAdditional))
				for _, add := range detail.BookingDetailsAdditional {
					additionalEntity := entity.BookingDetailAdditional{
						ID:                   add.ID,
						RoomTypeAdditionalID: add.RoomTypeAdditionalID,
						Category:             add.Category,
						Price:                add.Price,
						Pax:                  add.Pax,
						IsRequired:           add.IsRequired,
						NameAdditional:       add.NameAdditional,
					}
					additionalEntities = append(additionalEntities, additionalEntity)
					if strings.TrimSpace(add.NameAdditional) != "" {
						names = append(names, add.NameAdditional)
					}
				}
				result[i].BookingDetails[i2].BookingDetailsAdditional = additionalEntities
				result[i].BookingDetails[i2].BookingDetailAdditionalName = names
			}

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
