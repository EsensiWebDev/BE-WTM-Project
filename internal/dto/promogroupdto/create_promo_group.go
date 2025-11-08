package promogroupdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type CreatePromoGroupRequest struct {
	Name string `json:"name"`
}

func (r *CreatePromoGroupRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.Name, validation.Required.Error("Promo group name is required"), utils.NotEmptyAfterTrim("Name")))
}
