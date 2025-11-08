package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"mime/multipart"
	"wtm-backend/pkg/utils"
)

// RegisterRequest represents the input for user registration
type RegisterRequest struct {
	FullName     string                `form:"full_name"`
	AgentCompany string                `form:"agent_company"`
	Email        string                `form:"email"`
	Phone        string                `form:"phone"`
	Username     string                `form:"username" `
	KakaoTalkID  string                `form:"kakao_talk_id"`
	Password     string                `form:"password"`
	PhotoSelfie  *multipart.FileHeader `form:"photo_selfie"`
	PhotoIDCard  *multipart.FileHeader `form:"photo_id_card"`
	Certificate  *multipart.FileHeader `form:"certificate"`
	NameCard     *multipart.FileHeader `form:"name_card"`
}

func (r *RegisterRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FullName, validation.Required.Error("Full name is required"), utils.NotEmptyAfterTrim("Full Name")),
		validation.Field(&r.Email, validation.Required, is.Email.Error("Invalid email format"), utils.NotEmptyAfterTrim("Email")),
		validation.Field(&r.Phone, validation.Required, is.E164.Error("Phone number must use country code")),
		validation.Field(&r.Username, validation.Required.Error("Username is required"), utils.NotEmptyAfterTrim("Username")),
		validation.Field(&r.KakaoTalkID, validation.Required.Error("KakaoTaklId is required"), utils.NotEmptyAfterTrim("KakaoTalkId")),
		validation.Field(&r.Password, validation.Required.Error("Password is required")),
		validation.Field(&r.PhotoSelfie, validation.Required.Error("File selfie is required")),
		validation.Field(&r.PhotoIDCard, validation.Required.Error("File Id card is required")),
		validation.Field(&r.NameCard, validation.Required.Error("Name card is required")),
	)
}
