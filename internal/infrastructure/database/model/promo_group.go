package model

import "gorm.io/gorm"

type PromoGroup struct {
	gorm.Model
	ExternalID ExternalID `gorm:"embedded"`
	Name       string     `json:"name"`

	Promos []Promo `gorm:"many2many:detail_promo_groups"`
}

func (b *PromoGroup) BeforeCreate(tx *gorm.DB) error {
	return b.ExternalID.BeforeCreate(tx)
}
