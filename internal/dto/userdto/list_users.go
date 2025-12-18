package userdto

import (
	"strings"
	"wtm-backend/internal/dto"
	"wtm-backend/pkg/constant"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ListUsersRequest struct {
	Role                  string `form:"role"`
	AgentCompanyID        *uint  `form:"agent_company_id"`
	StatusID              uint   `form:"status_id"`
	dto.PaginationRequest `json:",inline"`
	Scope                 string `form:"scope"`
}

func (r *ListUsersRequest) Validate() error {
	r.Role = strings.ToLower(strings.TrimSpace(r.Role)) // normalisasi

	return validation.ValidateStruct(r,
		validation.Field(&r.Role, validation.Required,
			validation.In("", constant.RoleAdmin, constant.RoleSupport, constant.RoleAgent, constant.RoleSuperAdmin).
				Error("Role must be one of: admin, support, agent, super_admin"),
		),
	)
}

type ListUsersResponse struct {
	Users []ListUserData `json:"users"`
	Total int64          `json:"total"`
}

type ListUserData struct {
	ID               uint   `json:"id"`
	ExternalID       string `json:"external_id"`
	Name             string `json:"name"`
	Email            string `json:"email,omitempty"`
	Username         string `json:"username,omitempty"`
	PhoneNumber      string `json:"phone_number,omitempty"`
	Status           string `json:"status,omitempty"`
	PromoGroupName   string `json:"promo_group_name,omitempty"`
	AgentCompanyName string `json:"agent_company_name,omitempty"`
	KakaoTalkID      string `json:"kakao_talk_id,omitempty"`
	PromoGroupID     *uint  `json:"promo_group_id,omitempty"`
	Photo            string `json:"photo,omitempty"`
	Certificate      string `json:"certificate,omitempty"`
	NameCard         string `json:"name_card,omitempty"`
	IdCard           string `json:"id_card,omitempty"`
	Currency         string `json:"currency"`
}
