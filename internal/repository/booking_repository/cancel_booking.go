package booking_repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) CancelBooking(ctx context.Context, agentID uint, bookingDetailID string) error {
	db := br.db.GetTx(ctx)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Validasi booking detail milik agent dan masih waiting approval
		var exists bool
		if err := tx.Model(&model.BookingDetail{}).
			Select("count(*) > 0").
			Joins("JOIN bookings ON bookings.id = booking_details.booking_id").
			Where("booking_details.sub_booking_id = ?", bookingDetailID).
			Where("bookings.agent_id = ?", agentID).
			Where("booking_details.status_booking_id = ?", constant.StatusBookingWaitingApprovalID).
			Find(&exists).Error; err != nil {
			logger.Error(ctx, "failed to validate booking detail id with agent id", err.Error())
			return err
		}

		if !exists {
			logger.Error(ctx, "booking detail id not found for the agent or not waiting approval")
			return nil // nothing to cancel
		}

		// Step 2: Update detail ke canceled
		if err := tx.Model(&model.BookingDetail{}).
			Where("sub_booking_id = ?", bookingDetailID).
			Update("status_booking_id", constant.StatusBookingCanceledID).Error; err != nil {
			logger.Error(ctx, "failed to update booking detail status to canceled", err.Error())
			return err
		}

		// Step 3: Ambil bookingID dari detail ini
		var bookingID uint
		if err := tx.Model(&model.BookingDetail{}).
			Select("booking_id").
			Where("sub_booking_id = ?", bookingDetailID).
			Scan(&bookingID).Error; err != nil {
			logger.Error(ctx, "failed to get booking id", err.Error())
			return err
		}

		// Step 4: Cek apakah masih ada waiting approval di booking ini
		var waitingCount int64
		if err := tx.Model(&model.BookingDetail{}).
			Where("booking_id = ? AND status_booking_id = ?", bookingID, constant.StatusBookingWaitingApprovalID).
			Count(&waitingCount).Error; err != nil {
			logger.Error(ctx, "failed to count waiting approval status", err.Error())
			return err
		}

		// Step 5: Kalau tidak ada waiting approval, tentukan status prioritas
		if waitingCount == 0 {
			orderExpr := fmt.Sprintf(`
                CASE 
                    WHEN status_booking_id = %d THEN 1
                    WHEN status_booking_id = %d THEN 2
                    WHEN status_booking_id = %d THEN 3
                    ELSE 4
                END`,
				constant.StatusBookingRejectedID,
				constant.StatusBookingConfirmedID,
				constant.StatusBookingCanceledID,
			)

			var detailStatus uint
			if err := tx.Model(&model.BookingDetail{}).
				Select("status_booking_id").
				Where("booking_id = ?", bookingID).
				Order(orderExpr).
				Limit(1).
				Scan(&detailStatus).Error; err != nil {
				logger.Error(ctx, "failed to get booking status", err.Error())
				return err
			}

			// Step 6: Update Booking sesuai prioritas
			if err := tx.Model(&model.Booking{}).
				Where("id = ?", bookingID).
				Updates(map[string]interface{}{
					"status_booking_id": detailStatus,
					"approved_at":       time.Now(),
				}).Error; err != nil {
				logger.Error(ctx, "failed to update booking with priority", err.Error())
				return err
			}
		}

		return nil
	})
}
