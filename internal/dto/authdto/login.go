package authdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/utils"
)

// LoginRequest represents login input
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string          `json:"token"`
	User  *entity.UserMin `json:"user"`
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required.Error("Username is required"), utils.NotEmptyAfterTrim("Username")),
		validation.Field(&r.Password, validation.Required.Error("Password is required"), utils.NotEmptyAfterTrim("Password")),
	)
}
