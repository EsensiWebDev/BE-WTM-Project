package booking_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) CreateBookingDetail(ctx context.Context, detail *entity.BookingDetail) ([]uint, error) {
	db := br.db.GetTx(ctx)

	var baseDetail model.BookingDetail
	if err := utils.CopyPatch(&baseDetail, detail); err != nil {
		logger.Error(ctx, "Failed to copy booking detail entity to model", err.Error())
		return nil, err
	}

	countTrx := baseDetail.Quantity
	var ids []uint
	for i := 0; i < countTrx; i++ {
		var bookingDetail model.BookingDetail
		if err := utils.CopyPatch(&bookingDetail, baseDetail); err != nil {
			logger.Error(ctx, "Failed to copy booking detail entity to model", err.Error())
			return nil, err
		}
		code, err := br.GenerateCode(ctx, "sub_booking_codes", "SBK")
		if err != nil {
			logger.Error(ctx, "failed to generate sub booking code", err.Error())
			return nil, err
		}
		if code == "" {
			logger.Error(ctx, "failed to generate sub booking code after 10 attempts")
			return nil, err
		}
		bookingDetail.SubBookingID = code
		bookingDetail.Quantity = 1
		bookingDetail.ID = 0
		if err := db.WithContext(ctx).Create(&bookingDetail).Error; err != nil {
			logger.Error(ctx, "Failed to create booking detail", err.Error())
			return nil, err
		}
		ids = append(ids, bookingDetail.ID)
	}

	return ids, nil
}
