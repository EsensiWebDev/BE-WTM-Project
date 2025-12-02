package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type PromoType struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Name       string     `json:"name"`
}

type Promo struct {
	gorm.Model
	ExternalID  ExternalID     `gorm:"embedded"`
	Name        string         `json:"name"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	PromoTypeID uint           `json:"promo_type_id" gorm:"index"`
	Code        string         `json:"code" gorm:"unique"`
	Detail      datatypes.JSON `gorm:"type:jsonb"` // JSON field for additional details
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active"`

	PromoType      PromoType       `json:"promo_type" gorm:"foreignkey:PromoTypeID"`
	PromoGroups    []PromoGroup    `gorm:"many2many:detail_promo_groups"`
	PromoRoomTypes []PromoRoomType `gorm:"foreignkey:PromoID"`
}

func (b *Promo) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}

type PromoRoomType struct {
	gorm.Model
	ExternalID  ExternalID `gorm:"embedded"`
	PromoID     uint       `json:"promo_id" gorm:"index"`
	RoomTypeID  uint       `json:"room_type_id" gorm:"index"`
	TotalNights int        `json:"total_nights"`

	Promo    Promo    `json:"promo" gorm:"foreignkey:PromoID"`
	RoomType RoomType `json:"room_type" gorm:"foreignkey:RoomTypeID"`
}

func (b *PromoRoomType) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
