package authdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/pkg/utils"
)

// LoginRequest represents login input
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  DataUser `json:"user"`
}

type DataUser struct {
	RoleID      uint     `json:"role_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	PhotoURL    string   `json:"photo_url"`
	FullName    string   `json:"full_name"`
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required.Error("Username is required"), utils.NotEmptyAfterTrim("Username")),
		validation.Field(&r.Password, validation.Required.Error("Password is required"), utils.NotEmptyAfterTrim("Password")),
	)
}
