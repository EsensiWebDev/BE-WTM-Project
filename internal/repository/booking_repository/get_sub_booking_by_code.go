package booking_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetSubBookingByCode(ctx context.Context, code string) (*entity.BookingDetail, error) {
	db := br.db.GetTx(ctx)

	var bookingDetail model.BookingDetail
	if err := db.
		Preload("Booking").
		Where("sub_booking_id = ?", code).
		First(&bookingDetail).Error; err != nil {
		logger.Error(ctx, "failed to get booking by Id", err.Error())
		return nil, err
	}

	var result entity.BookingDetail
	if err := utils.CopyStrict(&result, &bookingDetail); err != nil {
		logger.Error(ctx, "failed to copy booking model to entity", err.Error())
		return nil, err
	}

	return &result, nil
}
