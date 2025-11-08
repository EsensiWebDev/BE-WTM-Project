package bannerdto

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateStatusBannerRequest struct {
	ID     uint `json:"id" validate:"required"`
	Status bool `json:"status" validate:"required"`
}

func (r *UpdateStatusBannerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required.Error("Banner Id is required")),
		validation.Field(&r.Status, validation.Required.Error("IsActive is required")),
	)
}
