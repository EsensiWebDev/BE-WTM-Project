package promogroupdto

import validation "github.com/go-ozzo/ozzo-validation"

type RemovePromoGroupMemberRequest struct {
	PromoGroupID uint `json:"promo_group_id"`
	MemberID     uint `json:"member_id"`
}

func (r *RemovePromoGroupMemberRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PromoGroupID, validation.Required.Error("Promo group Id is required")),
		validation.Field(&r.MemberID, validation.Required.Error("Member Id is required")),
	)
}
