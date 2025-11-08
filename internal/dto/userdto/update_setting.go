package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

type UpdateSettingRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (r *UpdateSettingRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required.Error("Username is required"), utils.NotEmptyAfterTrim("Username")),
		validation.Field(&r.OldPassword, validation.Required.Error("Old password is required"), utils.NotEmptyAfterTrim("Old Password")),
	)
}
