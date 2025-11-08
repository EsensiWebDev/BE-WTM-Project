package bookingdto

import "wtm-backend/internal/dto"

type ListBookingLogRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListBookingLogResponse struct {
	Data  []BookingLog `json:"data"`
	Total int64        `json:"total"`
}

type BookingLog struct {
	BookingCode   string `json:"booking_code"`
	ConfirmDate   string `json:"confirm_date"`
	AgentName     string `json:"agent_name"`
	BookingStatus string `json:"booking_status"`
	PaymentStatus string `json:"payment_status"`
	CheckInDate   string `json:"check_in_date"`
	CheckOutDate  string `json:"check_out_date"`
	HotelName     string `json:"hotel_name"`
	RoomTypeName  string `json:"room_type_name"`
	RoomNights    int    `json:"room_nights"`
	Capacity      string `json:"capacity"`
}
