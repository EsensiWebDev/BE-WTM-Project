package bookingdto

import "wtm-backend/internal/dto"

type ListBookingIDsRequest struct {
	dto.PaginationRequest `json:",inline"`
}
type ListBookingIDsResponse struct {
	BookingIDs []string `json:"booking_ids"`
	Total      int64    `json:"total"`
}
