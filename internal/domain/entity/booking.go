package entity

import (
	"time"
)

type Booking struct {
	ID              uint
	BookingCode     string
	AgentID         uint
	StatusBookingID uint
	StatusPaymentID uint
	BookingDetails  []BookingDetail

	Guests           []string       // Keep for backward compatibility
	BookingGuests    []BookingGuest // Full guest details
	BookingStatus    string
	PaymentStatus    string
	AgentName        string
	AgentCompanyName string
	AgentEmail       string
	AgentPhoneNumber string
	PromoGroupAgent  string
}

type BookingGuest struct {
	ID        uint
	BookingID uint
	Name      string
	Honorific string // e.g., "Mr", "Mrs", "Miss", "Ms"
	Category  string // "Adult" or "Child"
	Age       *int   // nullable, required when category="Child"
}

type BookingDetail struct {
	ID                          uint
	SubBookingID                string
	BookingID                   uint
	RoomPriceID                 uint
	CheckInDate                 time.Time
	CheckOutDate                time.Time
	Quantity                    int
	PromoID                     *uint
	DetailPromos                DetailPromo
	DetailRooms                 DetailRoom
	Price                       float64
	Currency                    string // Snapshot of currency at booking time
	Guest                       string
	OtherPreferences            string
	BedType                     string   // Selected bed type (singular)
	BedTypeNames                []string // Available bed types for selection
	AdditionalNotes             string   // Optional notes from agent to admin (max 500 characters)
	AdminNotes                  string   // Optional notes from admin to agent (max 500 characters)
	BookingDetailsAdditional    []BookingDetailAdditional
	RoomPrice                   RoomPrice
	StatusBookingID             uint
	StatusPaymentID             uint
	BookingDetailAdditionalName []string
	BookingStatus               string
	PaymentStatus               string
	ApprovedAt                  time.Time
	Booking                     Booking
	Invoice                     *Invoice
	ReceiptUrl                  string
	Promo                       *Promo
}

type DetailPromo struct {
	Name            string             `json:"name,omitempty"`
	PromoCode       string             `json:"promo_code,omitempty"`
	Type            string             `json:"type,omitempty"`
	Description     string             `json:"description,omitempty"`
	PromoTypeID     uint               `json:"promo_type_id,omitempty"`
	DiscountPercent float64            `json:"discount_percent,omitempty"`
	FixedPrice      float64            `json:"fixed_price,omitempty"`
	Prices          map[string]float64 `json:"prices,omitempty"` // Multi-currency prices
	UpgradedToID    uint               `json:"upgraded_to_id,omitempty"`
	BenefitNote     string             `json:"benefit_note,omitempty"`
}

type DetailRoom struct {
	HotelName     string `json:"hotel_name,omitempty"`
	RoomTypeName  string `json:"room_type_name,omitempty"`
	CancelledDate string `json:"cancelled_period,omitempty"`
	Capacity      int    `json:"capacity,omitempty"`
	IsAPI         bool   `json:"is_api,omitempty"`
}

type BookingDetailAdditional struct {
	ID                   uint
	BookingDetailIDs     []uint
	RoomTypeAdditionalID uint
	Category             string   // "price" or "pax"
	Price                *float64 // nullable, used when category="price"
	Pax                  *int     // nullable, used when category="pax"
	IsRequired           bool
	NameAdditional       string
}

// ReportAgentBooking represents a summary report of bookings made by agents
type ReportAgentBooking struct {
	AgentID          uint   `json:"agent_id"`
	AgentName        string `json:"agent_name"`
	AgentCompany     string `json:"agent_company"`
	HotelID          uint   `json:"hotel_id"`
	HotelName        string `json:"hotel_name"`
	ConfirmedBooking int64  `json:"confirmed_booking"`
	CancelledBooking int64  `json:"cancelled_booking"`
	RejectedBooking  int64  `json:"rejected_booking"`
}

// ReportAgentDetail represents detailed booking information for agents
type ReportAgentDetail struct {
	GuestName     string `json:"guest_name"`
	RoomType      string `json:"room_type"`
	DateIn        string `json:"date_in"`
	DateOut       string `json:"date_out"`
	Capacity      string `json:"capacity"`
	Additional    string `json:"additional"`
	StatusBooking string `json:"status_booking"`
}

type MonthlyBookingSummary struct {
	Month            string `json:"month"` // "2023-12"
	ConfirmedBooking int64  `json:"confirmed_booking"`
	CancelledBooking int64  `json:"cancelled_booking"`
	RejectedBooking  int64  `json:"rejected_booking"`
}

type MonthlyNewAgentSummary struct {
	Month    string `json:"month"` // "2023-12"
	NewAgent int64  `json:"new_agent"`
}

type ReportForGraph struct {
	DateTime *time.Time `json:"date_time,omitempty"`
	Date     string     `json:"date"`
	Count    int64      `json:"count"`
}

type StatusBooking struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

type StatusPayment struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}

type Invoice struct {
	BookingDetailID uint          `json:"booking_detail_id"`
	InvoiceCode     string        `json:"invoice_code"`
	DetailInvoice   DetailInvoice `json:"detail_invoice"`
	CreatedAt       time.Time     `json:"created_at"`
	BookingDetail   BookingDetail `json:"booking_detail"`
}

type DetailInvoice struct {
	CompanyAgent       string               `json:"company_agent"`
	Agent              string               `json:"agent"`
	Email              string               `json:"email"`
	Hotel              string               `json:"hotel"`
	Guest              string               `json:"guest"`
	CheckIn            string               `json:"check_in"`
	CheckOut           string               `json:"check_out"`
	SubBookingID       string               `json:"sub_booking_id"`
	BedType            string               `json:"bed_type,omitempty"`         // Selected bed type (e.g., "Kid Ogre Size")
	AdditionalNotes    string               `json:"additional_notes,omitempty"` // Optional notes for admin/agent only
	DescriptionInvoice []DescriptionInvoice `json:"description_invoice"`
	Promo              DetailPromo          `json:"promo"`
	Description        string               `json:"description"`
	TotalPrice         float64              `json:"total_price"`
	Currency           string               `json:"currency,omitempty"` // Currency code for the invoice (e.g. "IDR", "USD")
}

type DescriptionInvoice struct {
	Description      string  `json:"description"`
	Quantity         int     `json:"quantity"`
	Unit             string  `json:"unit"`
	Price            float64 `json:"price"`
	Total            float64 `json:"total"`
	TotalBeforePromo float64 `json:"total_before_promo,omitempty"`
	Category         string  `json:"category,omitempty"`    // "price" or "pax" - only for additional services
	Pax              *int    `json:"pax,omitempty"`         // nullable, used when category="pax"
	IsRequired       bool    `json:"is_required,omitempty"` // only for additional services
}
