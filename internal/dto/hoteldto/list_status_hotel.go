package hoteldto

import "wtm-backend/internal/domain/entity"

type ListStatusHotelResponse struct {
	StatusHotel []entity.StatusHotel `json:"status_hotel"`
}
