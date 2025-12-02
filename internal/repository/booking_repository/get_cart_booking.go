package booking_repository

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
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
		Preload("BookingDetails", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("BookingDetails.BookingDetailsAdditional").
		Preload("BookingDetails.RoomPrice.RoomType.RoomTypeAdditionals").
		Preload("BookingDetails.RoomPrice.RoomType.RoomTypeAdditionals.RoomAdditional").
		Preload("BookingGuests").
		Preload("BookingDetails.RoomPrice").
		Preload("BookingDetails.RoomPrice.RoomType").
		Preload("BookingDetails.RoomPrice.RoomType.Hotel").
		Preload("BookingDetails.Promo").
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
	for i, detail := range booking.BookingDetails {
		for _, additional := range detail.RoomPrice.RoomType.RoomTypeAdditionals {
			for i3, detailAdditional := range detail.BookingDetailsAdditional {
				if detailAdditional.RoomTypeAdditionalID == additional.ID {
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].NameAdditional = additional.RoomAdditional.Name
					bookingEntity.BookingDetails[i].BookingDetailsAdditional[i3].Price = additional.Price
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
		}
	}
	bookingEntity.Guests = guests

	return &bookingEntity, nil
}
