package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateBookingDetailGuest(ctx context.Context, detailID uint, guest string) error {
	db := br.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("id = ?", detailID).
		Update("guest", guest).Error

	if err != nil {
		logger.Error(ctx, "failed to update booking detail", "error", err)
		return fmt.Errorf("failed to update booking detail: %s", err.Error())
	}

	return nil
}
