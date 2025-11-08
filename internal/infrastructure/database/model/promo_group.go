package model

import "gorm.io/gorm"

type PromoGroup struct {
	gorm.Model
	Name string `json:"name"`

	Promos []Promo `gorm:"many2many:detail_promo_groups"`
}
