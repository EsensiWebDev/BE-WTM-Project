package booking_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"gorm.io/gorm"
)

func (br *BookingRepository) GetCartBooking(ctx context.Context, agentID uint) (*entity.Booking, error) {
	db := br.db.GetTx(ctx)

	var booking model.Booking
	if err := db.WithContext(ctx).
		Where("agent_id = ? AND status_booking_id = ?", agentID, constant.StatusBookingInCartID).
		Preload("BookingDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("BookingDetails.BookingDetailsAdditional").
		Preload("BookingDetails.RoomPrice.RoomType.RoomTypeAdditionals").
		Preload("BookingDetails.RoomPrice.RoomType.RoomTypeAdditionals.RoomAdditional").
		Preload("BookingDetails.RoomPrice.RoomType.BedTypes").
		Preload("BookingGuests").
		Preload("BookingDetails.RoomPrice").
		Preload("BookingDetails.RoomPrice.RoomType").
		Preload("BookingDetails.RoomPrice.RoomType.Hotel").
		Preload("BookingDetails.Promo").
		Preload("BookingDetails.Promo.PromoType").
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
	var bookingGuests []entity.BookingGuest
	for _, guest := range booking.BookingGuests {
		guests = append(guests, guest.Name)
		bookingGuests = append(bookingGuests, entity.BookingGuest{
			ID:        guest.ID,
			BookingID: guest.BookingID,
			Name:      guest.Name,
			Honorific: guest.Honorific,
			Category:  guest.Category,
			Age:       guest.Age,
		})
	}
	bookingEntity.BookingGuests = bookingGuests
	for i, detail := range booking.BookingDetails {
		// Map selected bed type from database
		bookingEntity.BookingDetails[i].BedType = detail.BedType

		// Map bed types from RoomType (available bed types for reference)
		var bedTypeNames []string
		for _, bedType := range detail.RoomPrice.RoomType.BedTypes {
			bedTypeNames = append(bedTypeNames, bedType.Name)
		}
		bookingEntity.BookingDetails[i].BedTypeNames = bedTypeNames

		// Map RoomAdditions for RoomType (needed for usecase to access Category, Pax, IsRequired)
		var roomAdditions []entity.CustomRoomAdditionalWithID
		for _, additional := range detail.RoomPrice.RoomType.RoomTypeAdditionals {
			roomAdditions = append(roomAdditions, entity.CustomRoomAdditionalWithID{
				ID:         additional.ID,
				Name:       additional.RoomAdditional.Name,
				Category:   additional.Category,
				Price:      additional.Price,
				Pax:        additional.Pax,
				IsRequired: additional.IsRequired,
			})
		}
		bookingEntity.BookingDetails[i].RoomPrice.RoomType.RoomAdditions = roomAdditions

		// Map additional services with new structure (Category, Price/Pax, IsRequired)
		// Now these fields are stored directly in BookingDetailAdditional, so we can read them from there
		for i3, detailAdditional := range detail.BookingDetailsAdditional {
			// Find matching RoomTypeAdditional to get NameAdditional if not already set
			for _, additional := range detail.RoomPrice.RoomType.RoomTypeAdditionals {
				if detailAdditional.RoomTypeAdditionalID == additional.ID {
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].NameAdditional = additional.RoomAdditional.Name
					// Map the stored fields from database model to entity
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].Category = detailAdditional.Category
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].Price = detailAdditional.Price
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].Pax = detailAdditional.Pax
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].IsRequired = detailAdditional.IsRequired
					break
				}
			}
		}
		if detail.Promo != nil {
			var detailPromo entity.PromoDetail
			if err := json.Unmarshal(detail.Promo.Detail, &detailPromo); err != nil {
				logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
			}
			bookingEntity.BookingDetails[i].Promo.Detail = detailPromo
			// Set PromoTypeName from the preloaded PromoType
			if detail.Promo.PromoType.Name != "" {
				bookingEntity.BookingDetails[i].Promo.PromoTypeName = detail.Promo.PromoType.Name
			}
		}
	}
	bookingEntity.Guests = guests

	return &bookingEntity, nil
}
