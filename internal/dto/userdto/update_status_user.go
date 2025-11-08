package userdto

import validation "github.com/go-ozzo/ozzo-validation"

type UpdateStatusUserRequest struct {
	ID       uint   `json:"id" form:"id"`
	IsActive bool   `json:"is_active" form:"is_active"`
	Reason   string `json:"reason" form:"reason"`
}

func (r *UpdateStatusUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required.Error("User Id is required")),
	)
}
