package bookingdto

import "wtm-backend/internal/dto"

type ListBookingsRequest struct {
	dto.PaginationRequest `json:",inline"`
	BookingStatusID       int `json:"booking_status_id" form:"booking_status_id"`
	PaymentStatusID       int `json:"payment_status_id" form:"payment_status_id"`
}

type ListBookingsResponse struct {
	Data  []DataBooking `json:"data"`
	Total int64         `json:"total"`
}

type DataBooking struct {
	BookingID     string          `json:"booking_id"`
	GuestName     []string        `json:"guest_name"`
	AgentName     string          `json:"agent_name"`
	AgentCompany  string          `json:"agent_company"`
	GroupPromo    string          `json:"group_promo"`
	BookingStatus string          `json:"booking_status"`
	PaymentStatus string          `json:"payment_status"`
	Receipts      []string        `json:"receipts"`
	Detail        []DetailBooking `json:"detail"`
}

type DetailBooking struct {
	GuestName     string      `json:"guest_name"`
	HotelName     string      `json:"hotel_name"`
	Additional    []string    `json:"additional"`
	SubBookingID  string      `json:"sub_booking_id"`
	BookingStatus string      `json:"booking_status"`
	PaymentStatus string      `json:"payment_status"`
	IsAPI         bool        `json:"is_api,omitempty"`
	CancelledDate string      `json:"cancelled_date,omitempty"`
	PromoCode     string      `json:"promo_code,omitempty"`
	Receipt       string      `json:"receipt_url,omitempty"`
	Invoice       DataInvoice `json:"invoice,omitempty"`
}
