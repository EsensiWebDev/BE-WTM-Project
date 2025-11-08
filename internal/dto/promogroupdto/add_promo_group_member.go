package promogroupdto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type AssignPromoGroupMemberRequest struct {
	PromoGroupID   uint `json:"promo_group_id"`
	MemberID       uint `json:"member_id"`
	AgentCompanyID uint `json:"agent_company_id"`
}

func (r *AssignPromoGroupMemberRequest) Validate() error {
	if err := validation.ValidateStruct(r,
		validation.Field(&r.PromoGroupID, validation.Required.Error("Promo group Id is required")),
	); err != nil {
		return err
	}

	if r.MemberID == 0 && r.AgentCompanyID == 0 {
		return validation.Errors{
			"member_id":        validation.NewInternalError(fmt.Errorf("either Member Id or Agent Company Id must be provided")),
			"agent_company_id": validation.NewInternalError(fmt.Errorf("either Member Id or Agent Company Id must be provided")),
		}
	}

	return nil
}
