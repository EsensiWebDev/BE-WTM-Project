package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UpdateStatus(ctx context.Context, req *hoteldto.UpdateStatusRequest) error {

	var statusId uint

	if req.Status {
		statusId = constant.StatusHotelApprovedID
	} else {
		statusId = constant.StatusHotelRejectedID
	}

	if err := hu.hotelRepo.UpdateStatus(ctx, req.HotelID, statusId); err != nil {
		logger.Error(ctx, " Error updating hotel status", err.Error())
		return err
	}

	return nil
}
