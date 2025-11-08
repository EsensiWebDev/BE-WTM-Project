package hotel_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListStatusHotel(ctx context.Context) (*hoteldto.ListStatusHotelResponse, error) {
	var statuses = constant.MapStatusHotel

	resp := &hoteldto.ListStatusHotelResponse{}

	if len(statuses) == 0 {
		logger.Error(ctx, "No hotel statuses found")
		return resp, nil
	}

	for id, status := range statuses {
		resp.StatusHotel = append(resp.StatusHotel, entity.StatusHotel{
			ID:     uint(id),
			Status: status,
		})
	}

	return resp, nil
}
