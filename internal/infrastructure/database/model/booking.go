package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StatusBooking struct {
	ID         uint       `gorm:"primaryKey"` // override default gorm.Model ID
	Status     string     `gorm:"type:text"`
	ExternalID ExternalID `gorm:"embedded"`
}

func (b *StatusBooking) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type StatusPayment struct {
	ID         uint       `gorm:"primaryKey"` // override default gorm.Model ID
	Status     string     `gorm:"type:text"`
	ExternalID ExternalID `gorm:"embedded"`
}

func (b *StatusPayment) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type Booking struct {
	gorm.Model
	ExternalID      ExternalID `gorm:"embedded"`
	AgentID         uint       `gorm:"index"`
	BookingCode     string     `gorm:"uniqueIndex;not null"`
	StatusBookingID uint       `gorm:"index"`
	StatusPaymentID uint       `gorm:"index"`

	StatusBooking  StatusBooking   `gorm:"foreignkey:StatusBookingID"`
	StatusPayment  StatusPayment   `gorm:"foreignkey:StatusPaymentID"`
	Agent          User            `gorm:"foreignkey:AgentID"`
	BookingDetails []BookingDetail `gorm:"foreignkey:BookingID"`
	BookingGuests  []BookingGuest  `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type BookingDetail struct {
	gorm.Model
	ExternalID   ExternalID `gorm:"embedded"`
	SubBookingID string     `gorm:"uniqueIndex;not null"`
	BookingID    uint       `gorm:"index"`
	RoomPriceID  uint       `gorm:"index"`
	CheckInDate  time.Time
	CheckOutDate time.Time
	ApprovedAt   time.Time
	RejectedAt   time.Time
	CancelledAt  time.Time
	Quantity     int
	ReceiptUrl   string `gorm:"type:text"`
	PaidAt       *time.Time

	// Promo snapshot
	PromoID     *uint          `gorm:"index"`      // nullable
	DetailPromo datatypes.JSON `gorm:"type:jsonb"` // snapshot of promo details
	DetailRoom  datatypes.JSON `gorm:"type:jsonb"` // snapshot of room details

	// Pricing
	Price    float64 `gorm:"type:float"`
	Currency string  `json:"currency" gorm:"type:varchar(3);default:'IDR'"` // Snapshot of currency at booking time

	// Guest per kamar
	Guest            string `gorm:"type:text"`
	BedType          string `gorm:"type:text"` // Selected bed type (e.g., "Kid Ogre Size")
	OtherPreferences string `gorm:"type:text"` // comma-separated list of selected other preferences
	AdditionalNotes  string `gorm:"type:text"` // Optional notes from agent to admin (max 500 characters)
	AdminNotes       string `gorm:"type:text"` // Optional notes from admin to agent (max 500 characters)

	// Status
	StatusBookingID uint `gorm:"index"`
	StatusPaymentID uint `gorm:"index"`

	Booking                  Booking                   `gorm:"foreignkey:BookingID"`
	Promo                    *Promo                    `gorm:"foreignkey:PromoID"`
	BookingDetailsAdditional []BookingDetailAdditional `gorm:"foreignkey:BookingDetailID"`
	RoomPrice                RoomPrice                 `gorm:"foreignkey:RoomPriceID"`

	StatusBooking StatusBooking `gorm:"foreignkey:StatusBookingID"`
	StatusPayment StatusPayment `gorm:"foreignkey:StatusPaymentID"`
	Invoice       *Invoice      `gorm:"foreignkey:BookingDetailID"`
}

func (b *BookingDetail) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type BookingGuest struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	BookingID  uint       `gorm:"index"`
	Name       string     `gorm:"type:text"`
	Honorific  string     `gorm:"type:varchar(10)"`                 // e.g., "Mr", "Mrs", "Miss", "Ms"
	Category   string     `gorm:"type:varchar(20);default:'Adult'"` // "Adult" or "Child"
	Age        *int       `gorm:"type:integer"`                     // nullable, required when category="Child"

	Booking Booking `gorm:"foreignkey:BookingID"`
}

func (b *BookingGuest) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type BookingDetailAdditional struct {
	gorm.Model
	ExternalID           ExternalID `gorm:"embedded"`
	BookingDetailID      uint       `gorm:"index"`
	RoomTypeAdditionalID uint       `gorm:"index"`
	Category             string     `gorm:"type:varchar(10);default:'price'"` // "price" or "pax"
	Price                *float64   `gorm:"type:decimal(10,2)"`               // nullable, used when category="price"
	Pax                  *int       `gorm:"type:integer"`                     // nullable, used when category="pax"
	IsRequired           bool       `gorm:"default:false"`
	NameAdditional       string     `gorm:"type:text"`
	Quantity             int

	BookingDetail BookingDetail `gorm:"foreignkey:BookingDetailID"`
}

func (b *BookingDetailAdditional) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type Invoice struct {
	gorm.Model
	ExternalID      ExternalID     `gorm:"embedded"`
	BookingDetailID uint           `gorm:"index;not null"`
	InvoiceCode     string         `gorm:"uniqueIndex;size:32;not null"`
	Detail          datatypes.JSON `gorm:"type:jsonb;not null"`

	BookingDetail BookingDetail `gorm:"foreignkey:BookingDetailID"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	return i.ExternalID.BeforeCreate(tx)
}
