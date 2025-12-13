package promodto

import validation "github.com/go-ozzo/ozzo-validation"

type SetStatusPromoRequest struct {
	PromoID  string `json:"promo_id" form:"promo_id"`
	IsActive bool   `json:"is_active" form:"is_active"`
}

func (r *SetStatusPromoRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PromoID, validation.Required.Error("Promo group Id is required")),
	)
}
