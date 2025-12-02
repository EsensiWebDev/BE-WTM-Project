package booking_repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateBookingReceipt(ctx context.Context, bookingDetailID []uint, receiptURL string) error {
	db := br.db.GetTx(ctx)

	if len(bookingDetailID) == 0 {
		logger.Error(ctx, "no bookingDetailID provided")
		return errors.New("no bookingDetailID provided")
	}

	updates := map[string]interface{}{
		"receipt_url": receiptURL,
		"paid_at":     gorm.Expr("NOW()"),
		//"status_payment_id": constant.StatusPaymentPaidID,
	}

	if err := db.Model(&model.BookingDetail{}).
		Where("id IN (?)", bookingDetailID).
		Updates(updates).Error; err != nil {
		logger.Error(ctx, "failed to update booking receipt", err.Error())
		return err
	}

	return nil
}
