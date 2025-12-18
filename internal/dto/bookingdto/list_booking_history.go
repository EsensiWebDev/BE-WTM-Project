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
	GuestName          string                     `json:"guest_name"`
	AgentName          string                     `json:"agent_name"`
	HotelName          string                     `json:"hotel_name"`
	RoomTypeName       string                     `json:"room_type_name,omitempty"`      // Room type selected
	IsBreakfast        bool                       `json:"is_breakfast"`                  // Whether breakfast is included
	BedType            string                     `json:"bed_type,omitempty"`            // Selected bed type
	RoomPrice          float64                    `json:"room_price"`                    // Room price per night (after promo if any)
	TotalPrice         float64                    `json:"total_price"`                   // Total price including room and additional services
	Currency           string                     `json:"currency,omitempty"`            // Currency code for prices
	CheckInDate        string                     `json:"check_in_date,omitempty"`       // Check-in date
	CheckOutDate       string                     `json:"check_out_date,omitempty"`      // Check-out date
	Additional         []string                   `json:"additional"`                    // Deprecated: use AdditionalServices for detailed info
	AdditionalServices []BookingHistoryAdditional `json:"additional_services,omitempty"` // Detailed additional services with price, category, pax, etc.
	OtherPreferences   []string                   `json:"other_preferences,omitempty"`
	SubBookingID       string                     `json:"sub_booking_id"`
	BookingStatus      string                     `json:"booking_status"`
	PaymentStatus      string                     `json:"payment_status"`
	CancellationDate   string                     `json:"cancellation_date"`
	Invoice            DataInvoice                `json:"invoice"`
	Receipt            string                     `json:"receipt"`
	AdditionalNotes    string                     `json:"additional_notes,omitempty"` // Notes from agent to admin
	AdminNotes         string                     `json:"admin_notes,omitempty"`      // Notes from admin to agent
}

type BookingHistoryAdditional struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`        // "price" or "pax"
	Price      *float64 `json:"price,omitempty"` // nullable, used when category="price"
	Pax        *int     `json:"pax,omitempty"`   // nullable, used when category="pax"
	IsRequired bool     `json:"is_required"`
}
