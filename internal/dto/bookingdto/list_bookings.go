package bookingdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

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
	BookingID     string             `json:"booking_id"`
	GuestName     []string           `json:"guest_name"`
	AgentName     string             `json:"agent_name"`
	AgentCompany  string             `json:"agent_company"`
	GroupPromo    string             `json:"group_promo"`
	PromoName     string             `json:"promo_name,omitempty"`
	DetailPromo   entity.DetailPromo `json:"detail_promo,omitempty"`
	BookingStatus string             `json:"booking_status"`
	PaymentStatus string             `json:"payment_status"`
	Receipts      []string           `json:"receipts"`
	Detail        []DetailBooking    `json:"detail"`
}

type DetailBooking struct {
	GuestName string `json:"guest_name"`
	HotelName string `json:"hotel_name"`
	// Room information (aligned with DetailBookingHistory / SubBookingDetail)
	RoomTypeName string  `json:"room_type_name,omitempty"` // Room type selected
	IsBreakfast  bool    `json:"is_breakfast"`             // Whether breakfast is included
	BedType      string  `json:"bed_type,omitempty"`       // Selected bed type
	RoomPrice    float64 `json:"room_price,omitempty"`     // Room price per night (after promo if any)
	TotalPrice   float64 `json:"total_price,omitempty"`    // Total price including room and services
	Currency     string  `json:"currency,omitempty"`       // Currency code for prices
	CheckInDate  string  `json:"check_in_date,omitempty"`  // Check-in date
	CheckOutDate string  `json:"check_out_date,omitempty"` // Check-out date

	// Additional services & preferences
	Additional         []string                   `json:"additional"`                    // Deprecated: use AdditionalServices for detailed info
	OtherPreferences   []string                   `json:"other_preferences,omitempty"`   // Simple text preferences
	AdditionalServices []BookingHistoryAdditional `json:"additional_services,omitempty"` // Detailed additional services with price, category, pax, etc.

	// Booking metadata
	SubBookingID  string `json:"sub_booking_id"`
	BookingStatus string `json:"booking_status"`
	PaymentStatus string `json:"payment_status"`
	IsAPI         bool   `json:"is_api,omitempty"`
	CancelledDate string `json:"cancelled_date,omitempty"`
	PromoCode     string `json:"promo_code,omitempty"`
	Receipt       string `json:"receipt_url,omitempty"`

	// Notes & invoice
	AdditionalNotes string      `json:"additional_notes,omitempty"` // Notes from agent to admin
	AdminNotes      string      `json:"admin_notes,omitempty"`      // Notes from admin to agent
	Invoice         DataInvoice `json:"invoice,omitempty"`
}
