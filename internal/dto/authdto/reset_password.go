package authdto

import validation "github.com/go-ozzo/ozzo-validation"

type ResetPasswordRequest struct {
	Token    string `json:"token" form:"token"`
	Password string `json:"password" form:"password"`
}

func (r *ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Token, validation.Required.Error("Token is required")),
		validation.Field(&r.Password, validation.Required.Error("Password is required")),
	)
}
