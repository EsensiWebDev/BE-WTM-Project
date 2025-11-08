package promogroupdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListPromoGroupMemberRequest struct {
	PromoGroupID          uint `json:"promo_group_id" form:"promo_group_id"`
	dto.PaginationRequest `json:",inline"`
}
type ListPromoGroupMemberResponse struct {
	PromoGroupMembers []ListPromoGroupMemberData `json:"promo_group_members"`
}

type ListPromoGroupMemberData struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	AgentCompany string `json:"agent_company"`
}

func (r *ListPromoGroupMemberRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.PromoGroupID, validation.Required.Error("Promo group Id is required")))
}
