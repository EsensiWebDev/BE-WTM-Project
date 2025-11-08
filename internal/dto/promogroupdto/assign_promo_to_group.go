package promogroupdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type AssignPromoToGroupRequest struct {
	PromoID      uint `json:"promo_id" form:"promo_id"`
	PromoGroupID uint `json:"promo_group_id" form:"promo_group_id"`
}

func (r *AssignPromoToGroupRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PromoID, validation.Required.Error("PromoID is required")),
		validation.Field(&r.PromoGroupID, validation.Required.Error("PromoGroupID is required")),
	)
}
