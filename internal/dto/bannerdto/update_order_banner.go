package bannerdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type UpdateOrderBannerRequest struct {
	ID    string `json:"id"`
	Order string `json:"order"`
}

func (r *UpdateOrderBannerRequest) Validate() error {

	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required.Error("Banner ID is required"), utils.NotEmptyAfterTrim("Banner ID")),
		validation.Field(&r.Order, validation.Required.Error("Order is required"), validation.In("up", "down").Error("Order must be up or down")),
	)
}
