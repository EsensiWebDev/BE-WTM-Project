package booking_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListStatusPayment(ctx context.Context) (*bookingdto.ListStatusPaymentResponse, error) {
	var statuses = constant.MapStatusPayment

	resp := &bookingdto.ListStatusPaymentResponse{}

	if len(statuses) == 0 {
		logger.Error(ctx, "No status payments found")
		return resp, nil
	}

	for _, id := range constant.StatusPaymentOrder {
		status := statuses[id]
		resp.Data = append(resp.Data, entity.StatusPayment{
			ID:     uint(id),
			Status: status,
		})
	}

	return resp, nil
}
