package bookingdto

import "wtm-backend/internal/domain/entity"

type ListStatusBookingResponse struct {
	Data []entity.StatusBooking `json:"bookings"`
}
