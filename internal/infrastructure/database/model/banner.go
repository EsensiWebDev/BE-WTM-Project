package model

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	ExternalID   ExternalID `gorm:"embedded"`
	Title        string     `json:"title"`
	ImageURL     string     `json:"image_url"`
	Description  string     `json:"description"`
	IsActive     bool       `json:"is_active"`
	DisplayOrder int        `json:"display_order"`
}

func (b *Banner) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
