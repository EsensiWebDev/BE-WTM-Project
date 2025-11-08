package authdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"wtm-backend/pkg/utils"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" form:"email"`
}

type ForgotPasswordResponse struct {
	ExpiresAt string `json:"expires_at"`
}

func (r *ForgotPasswordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email.Error("Invalid email format"), utils.NotEmptyAfterTrim("Email")),
	)
}
