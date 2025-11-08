package userdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListUsersByRoleRequest struct {
	Role                  string `json:"role"`
	dto.PaginationRequest `json:",inline"`
}

type ListUsersByRoleResponse struct {
	User []ListUsersByRoleData `json:"user"`
}

type ListUsersByRoleData struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	Status           string `json:"status"`
	PromoGroupID     uint   `json:"promo_group_id,omitempty"`
	PromoGroupName   string `json:"promo_group_name,omitempty"`
	AgentCompanyName string `json:"agent_company_name,omitempty"`
	KakaoTalkID      string `json:"kakao_talk_id,omitempty"`
}

func (r *ListUsersByRoleRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.Role, validation.Required, validation.In("admin", "support", "agent", "super_admin").Error("Role must be one of: admin, support, agent, super_admin")))
}
