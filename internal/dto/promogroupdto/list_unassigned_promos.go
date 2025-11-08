package promogroupdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListUnassignedPromosRequest struct {
	PromoGroupID          uint `json:"promo_group_id" form:"promo_group_id"`
	dto.PaginationRequest `json:",inline"`
}

func (r *ListUnassignedPromosRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PromoGroupID, validation.Required.Error("PromoGroupID is required")),
	)
}

type ListUnassignedPromosResponse struct {
	Promos []ListUnassignedPromoData
	Total  int64
}

type ListUnassignedPromoData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
