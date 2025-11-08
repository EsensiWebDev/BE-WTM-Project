package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"wtm-backend/pkg/utils"
)

type UpdateProfileRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	KakaoTalkID string `json:"kakao_talk_id"`
}

type UpdateProfileResponse struct {
}

func (r *UpdateProfileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FullName, validation.Required.Error("Full name is required"), utils.NotEmptyAfterTrim("Full Name")),
		validation.Field(&r.Email, validation.Required, is.Email.Error("Invalid email format"), utils.NotEmptyAfterTrim("Email")),
		validation.Field(&r.Phone, validation.Required, is.E164.Error("Phone number must use country code")),
	)
}
