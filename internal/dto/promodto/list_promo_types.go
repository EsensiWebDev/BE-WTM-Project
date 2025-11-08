package promodto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListPromoTypesRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListPromoTypesResponse struct {
	PromoTypes []entity.PromoType `json:"promo_types"`
}
