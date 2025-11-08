package booking_repository

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) UpdateBookingDetailStatus(ctx context.Context, bookingDetailIDs []uint, statusID uint) ([]entity.BookingDetail, error) {
	db := br.db.GetTx(ctx)

	var updatedDetails []model.BookingDetail
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Update BookingDetail
		if err := tx.Model(&model.BookingDetail{}).
			Where("id IN ?", bookingDetailIDs).
			Updates(map[string]interface{}{
				"status_booking_id": statusID,
				"approved_at":       time.Now(),
			}).Error; err != nil {
			logger.Error(ctx, "failed to update booking detail status", err.Error())
			return err
		}

		// Step 2: Ambil BookingID yang terkait
		var bookingIDs []uint
		if err := tx.Model(&model.BookingDetail{}).
			Where("id IN ?", bookingDetailIDs).
			Distinct().
			Pluck("booking_id", &bookingIDs).Error; err != nil {
			logger.Error(ctx, "failed to fetch booking IDs", err.Error())
			return err
		}

		// Step 3: Update Booking
		if err := tx.Model(&model.Booking{}).
			Where("id IN ?", bookingIDs).
			Updates(map[string]interface{}{
				"status_booking_id": statusID,
				"approved_at":       time.Now(),
			}).Error; err != nil {
			logger.Error(ctx, "failed to update booking status", err.Error())
			return err
		}

		// Step 4: Fetch updated BookingDetails with preload
		if err := tx.Preload("Booking").
			Preload("Booking.Agent").
			Preload("RoomType").
			Preload("BookingDetailsAdditional").
			Where("id IN ?", bookingDetailIDs).
			Find(&updatedDetails).Error; err != nil {
			logger.Error(ctx, "failed to fetch updated booking details", err.Error())
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(ctx, "failed to update status booking", err.Error())
		return nil, err
	}

	var result []entity.BookingDetail
	if err := utils.CopyStrict(&result, &updatedDetails); err != nil {
		logger.Error(ctx, "failed to copy updated booking details", err.Error())
		return nil, err
	}

	for i, detail := range updatedDetails {
		var detailRoom entity.DetailRoom
		if err := json.Unmarshal(detail.DetailRoom, &detailRoom); err != nil {
			logger.Error(ctx, "Error marshalling room detail to JSON", err.Error())
		}
		result[i].DetailRooms = detailRoom

		if len(detail.BookingDetailsAdditional) > 0 {
			result[i].BookingDetailAdditionalName = make([]string, 0, len(detail.BookingDetailsAdditional))
			for _, additional := range detail.BookingDetailsAdditional {
				result[i].BookingDetailAdditionalName = append(result[i].BookingDetailAdditionalName, additional.NameAdditional)
			}
		}
		if detail.Booking.ID != 0 {
			result[i].Booking.AgentName = detail.Booking.Agent.FullName
			result[i].Booking.AgentEmail = detail.Booking.Agent.Email
		}
	}

	return result, nil
}
