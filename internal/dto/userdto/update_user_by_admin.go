package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"wtm-backend/pkg/utils"
)

type UpdateUserByAdminRequest struct {
	UserID                   uint `json:"user_id" form:"user_id"`
	CreateUserByAdminRequest `json:",inline"`
	IsActive                 bool `json:"is_active" form:"is_active"`
}

func (r *UpdateUserByAdminRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required.Error("User Id is required")),
		validation.Field(&r.FullName, validation.Required.Error("Full name is required"), utils.NotEmptyAfterTrim("Full Name")),
		validation.Field(&r.Email, validation.Required, is.Email.Error("Invalid email format"), utils.NotEmptyAfterTrim("Email")),
		validation.Field(&r.Phone, validation.Required, is.E164.Error("Phone number must use country code")),
	)
}
