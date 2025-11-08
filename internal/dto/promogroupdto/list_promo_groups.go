package promogroupdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListPromoGroupRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListPromoGroupResponse struct {
	PromoGroups []entity.PromoGroup `json:"promo_groups"`
}
