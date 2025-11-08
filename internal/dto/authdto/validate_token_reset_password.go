package authdto

import validation "github.com/go-ozzo/ozzo-validation"

type ValidateTokenResetPasswordRequest struct {
	Token string `json:"token" form:"token"`
}

type ValidateTokenResetPasswordResponse struct {
	Email string `json:"email"`
}

func (r *ValidateTokenResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Token, validation.Required.Error("Token is required")),
	)
}
