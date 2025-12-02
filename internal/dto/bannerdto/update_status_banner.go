package bannerdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type UpdateStatusBannerRequest struct {
	ID     string `json:"id"`
	Status bool   `json:"status"`
}

func (r *UpdateStatusBannerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required.Error("Banner ID is required"), utils.NotEmptyAfterTrim("Banner ID")),
	)
}
