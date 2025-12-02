package bookingdto

import "wtm-backend/internal/dto"

type ListBookingLogRequest struct {
	dto.PaginationRequest `json:",inline"`
	BookingStatusID       int    `json:"booking_status_id" form:"booking_status_id"`
	PaymentStatusID       int    `json:"payment_status_id" form:"payment_status_id"`
	ConfirmDateFrom       string `json:"confirm_date_from" form:"confirm_date_from"`
	ConfirmDateTo         string `json:"confirm_date_to" form:"confirm_date_to"`
	CheckInDateFrom       string `json:"check_in_date_from" form:"check_in_date_from"`
	CheckInDateTo         string `json:"check_in_date_to" form:"check_in_date_to"`
	CheckOutDateFrom      string `json:"check_out_date_from" form:"check_out_date_from"`
	CheckOutDateTo        string `json:"check_out_date_to" form:"check_out_date_to"`
}

type ListBookingLogResponse struct {
	Data  []BookingLog `json:"data"`
	Total int64        `json:"total"`
}

type BookingLog struct {
	SubBookingID  string `json:"sub_booking_id"`
	BookingID     string `json:"booking_id"`
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
