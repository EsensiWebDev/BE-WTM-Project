package booking_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListStatusBooking(ctx context.Context) (*bookingdto.ListStatusBookingResponse, error) {
	var statuses = constant.MapStatusBooking

	resp := &bookingdto.ListStatusBookingResponse{}

	if len(statuses) == 0 {
		logger.Error(ctx, "No status bookings found")
		return resp, nil
	}

	for id, status := range statuses {
		resp.Data = append(resp.Data, entity.StatusBooking{
			ID:     uint(id),
			Status: status,
		})
	}

	return resp, nil
}
