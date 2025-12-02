package booking_repository

import (
	"context"
	"gorm.io/gorm"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateBookingDetailStatusPayment(ctx context.Context, bookingDetailIDs []uint, statusID uint) error {
	db := br.db.GetTx(ctx)
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Update BookingDetail yang status payment-nya berubah
		if err := tx.Model(&model.BookingDetail{}).
			Where("id IN ?", bookingDetailIDs).
			Updates(map[string]interface{}{
				"status_payment_id": statusID,
			}).Error; err != nil {
			return err
		}

		// Step 2: Ambil bookingID (cukup satu, karena pasti sama)
		var bookingID uint
		if err := tx.Model(&model.BookingDetail{}).
			Select("booking_id").
			Where("id = ?", bookingDetailIDs[0]).
			Scan(&bookingID).Error; err != nil {
			return err
		}

		// Step 3: Cek apakah masih ada detail unpaid
		var unpaidCount int64
		if err := tx.Model(&model.BookingDetail{}).
			Where("booking_id = ? AND status_payment_id = ?", bookingID, constant.StatusPaymentUnpaidID).
			Count(&unpaidCount).Error; err != nil {
			return err
		}

		// Step 4: Tentukan status parent
		var parentStatus uint
		if unpaidCount > 0 {
			parentStatus = constant.StatusPaymentUnpaidID
		} else {
			parentStatus = statusID
		}

		// Step 5: Update Booking parent
		if err := tx.Model(&model.Booking{}).
			Where("id = ?", bookingID).
			Updates(map[string]interface{}{
				"status_payment_id": parentStatus,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(ctx, "failed to update status booking", err.Error())
		return err
	}

	return nil
}
