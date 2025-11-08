package promogroupdto

import validation "github.com/go-ozzo/ozzo-validation"

type RemovePromoFromGroupRequest struct {
	PromoGroupID uint `json:"promo_group_id" form:"promo_group_id"`
	PromoID      uint `json:"promo_id" form:"promo_id"`
}

func (r *RemovePromoFromGroupRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PromoGroupID, validation.Required.Error("PromoGroupID is required")),
		validation.Field(&r.PromoID, validation.Required.Error("PromoID is required")),
	)
}
