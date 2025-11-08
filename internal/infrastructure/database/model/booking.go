package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type StatusBooking struct {
	ID     uint   `gorm:"primaryKey"` // override default gorm.Model ID
	Status string `gorm:"type:text"`
}

type StatusPayment struct {
	ID     uint   `gorm:"primaryKey"` // override default gorm.Model ID
	Status string `gorm:"type:text"`
}
type Booking struct {
	gorm.Model
	AgentID         uint   `gorm:"index"`
	BookingCode     string `gorm:"uniqueIndex;not null"`
	StatusBookingID uint   `gorm:"index"`
	StatusPaymentID uint   `gorm:"index"`
	ApprovedAt      time.Time

	StatusBooking  StatusBooking   `gorm:"foreignkey:StatusBookingID"`
	StatusPayment  StatusPayment   `gorm:"foreignkey:StatusPaymentID"`
	Agent          User            `gorm:"foreignkey:AgentID"`
	BookingDetails []BookingDetail `gorm:"foreignkey:BookingID"`
	BookingGuests  []BookingGuest  `gorm:"foreignkey:BookingID"`
}

type BookingDetail struct {
	gorm.Model
	SubBookingID string `gorm:"uniqueIndex;not null"`
	BookingID    uint   `gorm:"index"`
	RoomTypeID   uint   `gorm:"index"`
	CheckInDate  time.Time
	CheckOutDate time.Time
	ApprovedAt   time.Time
	Quantity     int

	// Promo snapshot
	PromoID     *uint          `gorm:"index"`      // nullable
	DetailPromo datatypes.JSON `gorm:"type:jsonb"` // snapshot of promo details
	DetailRoom  datatypes.JSON `gorm:"type:jsonb"` // snapshot of room details

	// Pricing
	Price float64 `gorm:"type:float"`

	// Guest per kamar
	Guest string `gorm:"type:text"`

	// Status
	StatusBookingID uint `gorm:"index"`
	StatusPaymentID uint `gorm:"index"`

	Booking                  Booking                   `gorm:"foreignkey:BookingID"`
	RoomType                 RoomType                  `gorm:"foreignkey:RoomTypeID"`
	Promo                    Promo                     `gorm:"foreignkey:PromoID"`
	BookingDetailsAdditional []BookingDetailAdditional `gorm:"foreignkey:BookingDetailID"`

	StatusBooking StatusBooking `gorm:"foreignkey:StatusBookingID"`
	StatusPayment StatusPayment `gorm:"foreignkey:StatusPaymentID"`
}

type BookingGuest struct {
	gorm.Model
	BookingID uint   `gorm:"index"`
	Name      string `gorm:"type:text"`

	Booking Booking `gorm:"foreignkey:BookingID"`
}

type BookingDetailAdditional struct {
	gorm.Model
	BookingDetailID      uint `gorm:"index"`
	RoomTypeAdditionalID uint `gorm:"index"`
	Price                float64
	NameAdditional       string `gorm:"type:text"`
	Quantity             int

	BookingDetail      BookingDetail      `gorm:"foreignkey:BookingDetailID"`
	RoomTypeAdditional RoomTypeAdditional `gorm:"foreignkey:RoomTypeAdditionalID"`
}

type Invoice struct {
	gorm.Model
	BookingDetailID uint           `gorm:"index;not null"`
	InvoiceCode     string         `gorm:"uniqueIndex;size:32;not null"`
	Detail          datatypes.JSON `gorm:"type:jsonb;not null"`

	BookingDetail BookingDetail `gorm:"foreignkey:BookingDetailID"`
}
