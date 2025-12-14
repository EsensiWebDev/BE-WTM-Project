package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RoomType struct {
	gorm.Model
	ExternalID       ExternalID     `gorm:"embedded"`
	HotelID          uint           `json:"hotel_id" gorm:"index"`
	Name             string         `json:"name"`
	IsSmokingAllowed bool           `json:"is_smoking_allowed"`
	MaxOccupancy     int            `json:"max_occupancy"`
	RoomSize         float64        `json:"room_size"` // in square meters
	Description      string         `json:"description"`
	Photos           pq.StringArray `json:"photos" gorm:"type:text[]"`
	TotalUnit        int            `json:"total_unit"`

	Hotel Hotel `gorm:"foreignKey:HotelID"`

	BedTypes   []BedType   `gorm:"many2many:bed_type_rooms"`
	RoomPrices []RoomPrice `gorm:"foreignKey:RoomTypeID"`

	RoomTypeAdditionals []RoomTypeAdditional `gorm:"foreignKey:RoomTypeID"`
	PromoRoomTypes      []PromoRoomType      `gorm:"foreignkey:RoomTypeID"`
}

func (b *RoomType) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type RoomPrice struct {
	gorm.Model
	ExternalID  ExternalID `gorm:"embedded"`
	RoomTypeID  uint       `json:"room_type_id_id" gorm:"index"`
	IsBreakfast bool       `json:"is_breakfast"`
	Pax         int        `json:"pax"`
	Price       float64    `json:"price"`
	IsShow      bool       `json:"is_show"`

	RoomType RoomType `gorm:"foreignKey:RoomTypeID"`
}

func (b *RoomPrice) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type RoomAdditional struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Name       string     `json:"name"`

	RoomTypeAdditionals []RoomTypeAdditional `gorm:"foreignKey:RoomAdditionalID"`
}

func (b *RoomAdditional) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type RoomTypeAdditional struct {
	gorm.Model
	ExternalID       ExternalID `gorm:"embedded"`
	RoomTypeID       uint       `json:"room_type_id" gorm:"index"`
	RoomAdditionalID uint       `json:"room_additional_id" gorm:"index"`
	Category         string     `json:"category" gorm:"type:varchar(10);default:'price'"` // "price" or "pax"
	Price            *float64   `json:"price" gorm:"type:decimal(10,2)"`                  // nullable, used when category="price"
	Pax              *int       `json:"pax"`                                              // nullable, used when category="pax"
	IsRequired       bool       `json:"is_required" gorm:"default:false"`

	RoomType       RoomType       `json:"room_type" gorm:"foreignkey:RoomTypeID"`
	RoomAdditional RoomAdditional `json:"room_additional" gorm:"foreignkey:RoomAdditionalID"`
}

func (b *RoomTypeAdditional) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type BedType struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Name       string     `json:"name"`

	RoomType []RoomType `json:"variants" gorm:"many2many:BedTypeRoom"`
}

func (b *BedType) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type RoomUnavailable struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	RoomTypeID uint       `json:"room_type_id" gorm:"index"`
	Date       *time.Time `json:"date"` // Date with format "2006-01-02"
	Reason     string     `json:"reason"`
	Source     string     `json:"source"`
	RoomType   RoomType   `gorm:"foreignkey:RoomTypeID"`
}

func (b *RoomUnavailable) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
