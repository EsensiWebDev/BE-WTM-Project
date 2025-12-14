package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) AddGuestsToCart(ctx context.Context, agentID uint, bookingID uint, guests []bookingdto.GuestInfo) error {
	db := br.db.GetTx(ctx)

	// Step 1: Validasi booking ID
	var exists bool
	if err := db.WithContext(ctx).
		Model(&model.Booking{}).
		Select("count(*) > 0").
		Where("id = ?", bookingID).
		Where("agent_id = ?", agentID).
		Find(&exists).Error; err != nil {
		logger.Error(ctx, "failed to validate booking ID", err.Error())
		return err
	}
	if !exists {
		return fmt.Errorf("booking ID %d not found", bookingID)
	}

	// Step 2: Buat slice BookingGuest dengan fields baru
	var guestModels []model.BookingGuest
	for _, guest := range guests {
		guestModels = append(guestModels, model.BookingGuest{
			BookingID: bookingID,
			Name:      guest.Name,
			Honorific: guest.Honorific,
			Category:  guest.Category,
			Age:       guest.Age,
		})
	}

	// Step 3: Insert batch
	if len(guestModels) > 0 {
		if err := db.WithContext(ctx).
			Create(&guestModels).Error; err != nil {
			logger.Error(ctx, "failed to insert booking guests", err.Error())
			return err
		}
	}

	return nil

}
