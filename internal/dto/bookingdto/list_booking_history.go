package bookingdto

import "wtm-backend/internal/dto"

type ListBookingHistoryRequest struct {
	dto.PaginationRequest `json:",inline"`
	SearchBy              string `json:"search_by" form:"search_by"`
	StatusBookingID       int    `json:"status_booking_id" form:"status_booking_id"`
	StatusPaymentID       int    `json:"status_payment_id" form:"status_payment_id"`
}

type ListBookingHistoryResponse struct {
	Data  []DataBookingHistory `json:"data"`
	Total int64                `json:"total"`
}

type DataBookingHistory struct {
	BookingID     uint                   `json:"booking_id"`
	GuestName     []string               `json:"guest_name"`
	BookingCode   string                 `json:"booking_code"`
	BookingStatus string                 `json:"booking_status"`
	PaymentStatus string                 `json:"payment_status"`
	Detail        []DetailBookingHistory `json:"detail"`
	Invoices      []DataInvoice          `json:"invoices"`
	Receipts      []string               `json:"receipts"`
}

type DetailBookingHistory struct {
	GuestName        string      `json:"guest_name"`
	AgentName        string      `json:"agent_name"`
	HotelName        string      `json:"hotel_name"`
	Additional       []string    `json:"additional"`
	SubBookingID     string      `json:"sub_booking_id"`
	BookingStatus    string      `json:"booking_status"`
	PaymentStatus    string      `json:"payment_status"`
	CancellationDate string      `json:"cancellation_date"`
	Invoice          DataInvoice `json:"invoice"`
	Receipt          string      `json:"receipt"`
}
