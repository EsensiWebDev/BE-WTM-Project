package bookingdto

import "wtm-backend/internal/dto"

type ListBookingsRequest struct {
	dto.PaginationRequest `json:",inline"`
	StatusBookingID       uint `json:"status_booking_id"`
}

type ListBookingsResponse struct {
	Data  []DataBooking `json:"data"`
	Total int64         `json:"total"`
}

type DataBooking struct {
	BookingID     uint            `json:"booking_id"`
	GuestName     []string        `json:"guest_name"`
	AgentName     string          `json:"agent_name"`
	AgentCompany  string          `json:"agent_company"`
	GroupPromo    string          `json:"group_promo"`
	BookingCode   string          `json:"booking_code"`
	BookingStatus string          `json:"booking_status"`
	PaymentStatus string          `json:"payment_status"`
	Detail        []DetailBooking `json:"detail"`
}

type DetailBooking struct {
	GuestName     string   `json:"guest_name"`
	HotelName     string   `json:"hotel_name"`
	Additional    []string `json:"additional"`
	SubBookingID  string   `json:"sub_booking_id"`
	BookingStatus string   `json:"booking_status"`
	PaymentStatus string   `json:"payment_status"`
	IsAPI         bool     `json:"is_api,omitempty"`
	CancelledDate string   `json:"cancelled_date,omitempty"`
	PromoID       *uint    `json:"promo_id,omitempty"`
	PromoCode     string   `json:"promo_code,omitempty"`
}
