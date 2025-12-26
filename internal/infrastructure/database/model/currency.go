package model

import (
	"gorm.io/gorm"
)

type Currency struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Code       string     `json:"code" gorm:"uniqueIndex;type:varchar(3);not null"` // ISO 4217: USD, IDR, EUR, etc.
	Name       string     `json:"name" gorm:"type:varchar(100);not null"`           // "US Dollar", "Indonesian Rupiah"
	Symbol     string     `json:"symbol" gorm:"type:varchar(10)"`                   // "$", "IDR", "€"
	IsActive   bool       `json:"is_active" gorm:"default:true"`
}

func (c *Currency) BeforeCreate(tx *gorm.DB) error {
	return c.ExternalID.BeforeCreate(tx)
}
