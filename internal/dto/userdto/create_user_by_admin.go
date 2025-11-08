package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"mime/multipart"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/utils"
)

// CreateUserByAdminRequest represents the request payload for creating a new user.
type CreateUserByAdminRequest struct {
	FullName     string                `json:"full_name" form:"full_name"`
	AgentCompany string                `json:"agent_company" form:"agent_company"`
	Email        string                `json:"email" form:"email"`
	Phone        string                `json:"phone" form:"phone"`
	PromoGroupID uint                  `json:"promo_group_id" form:"promo_group_id"`
	Role         string                `json:"role" form:"role"` // e.g., "admin", "suppor", "agent", "super_admin"
	KakaoTalkID  string                `form:"kakao_talk_id" json:"kakao_talk_id"`
	PhotoSelfie  *multipart.FileHeader `form:"photo_selfie" json:"photo_selfie"`
	PhotoIDCard  *multipart.FileHeader `form:"photo_id_card" json:"photo_id_card"`
	Certificate  *multipart.FileHeader `form:"certificate" json:"certificate"`
	NameCard     *multipart.FileHeader `form:"name_card" json:"name_card"`
}

func (r *CreateUserByAdminRequest) Validate() error {
	if r.Role == constant.RoleAgent {
		return validation.ValidateStruct(r,
			validation.Field(&r.AgentCompany, validation.Required.Error("Agent Company is required for agent role"), utils.NotEmptyAfterTrim("Agent Company")),
		)
	}
	return validation.ValidateStruct(r,
		validation.Field(&r.FullName, validation.Required.Error("Full name is required"), utils.NotEmptyAfterTrim("Full Name")),
		validation.Field(&r.Email, validation.Required, is.Email.Error("Invalid email format"), utils.NotEmptyAfterTrim("Email")),
		validation.Field(&r.Phone, validation.Required, is.E164.Error("Phone number must use country code")),
		validation.Field(&r.Role, validation.Required, validation.In("admin", "support", "agent", "super_admin").Error("Role must be one of: admin, support, agent, super_admin")),
	)
}
