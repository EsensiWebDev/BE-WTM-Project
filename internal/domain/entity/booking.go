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

	Guests           []string
	BookingStatus    string
	PaymentStatus    string
	AgentName        string
	AgentCompanyName string
	AgentEmail       string
	AgentPhoneNumber string
}

type BookingDetail struct {
	ID                          uint
	SubBookingID                string
	BookingID                   uint
	RoomTypeID                  uint
	CheckInDate                 time.Time
	CheckOutDate                time.Time
	Quantity                    int
	PromoID                     *uint
	DetailPromos                DetailPromo
	DetailRooms                 DetailRoom
	Price                       float64
	Guest                       string
	BookingDetailAdditional     []BookingDetailAdditional
	RoomPrice                   RoomPrice
	StatusBookingID             uint
	StatusPaymentID             uint
	BookingDetailAdditionalName []string
	BookingStatus               string
	PaymentStatus               string
	UpdatedAt                   time.Time
	Booking                     Booking
	RoomType                    RoomType
}

type DetailPromo struct {
	PromoCode       string  `json:"promo_code,omitempty"`
	Type            string  `json:"type,omitempty"`
	DiscountPercent float64 `json:"discount_percent,omitempty"`
	FixedPrice      float64 `json:"fixed_price,omitempty"`
	UpgradedToID    uint    `json:"upgraded_to_id,omitempty"`
	BenefitNote     string  `json:"benefit_note,omitempty"`
}

type DetailRoom struct {
	HotelName     string `json:"hotel_name,omitempty"`
	HotelPhoto    string `json:"hotel_photo,omitempty"`
	HotelRating   int    `json:"hotel_rating,omitempty"`
	RoomTypeName  string `json:"room_type_name,omitempty"`
	IsBreakfast   bool   `json:"is_breakfast,omitempty"`
	CancelledDate string `json:"cancelled_period,omitempty"`
	Capacity      string `json:"capacity,omitempty"`
	IsAPI         bool   `json:"is_api,omitempty"`
}

type BookingDetailAdditional struct {
	ID                   uint
	BookingDetailIDs     []uint
	RoomTypeAdditionalID uint
	Price                float64
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
	Month               string `json:"month"` // "2023-12"
	ConfirmedBooking    int64  `json:"confirmed_booking"`
	CancellationBooking int64  `json:"cancellation_booking"`
}

type MonthlyNewAgentSummary struct {
	Month    string `json:"month"` // "2023-12"
	NewAgent int64  `json:"new_agent"`
}

type ReportForGraph struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type StatusBooking struct {
	ID     uint
	Status string
}

type StatusPayment struct {
	ID     uint
	Status string
}
