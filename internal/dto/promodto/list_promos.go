package promodto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListPromosRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListPromosResponse struct {
	Promos []PromoResponse `json:"promos"`
}

type PromoResponse struct {
	ID               uint               `json:"id"`
	PromoCode        string             `json:"promo_code"`
	PromoName        string             `json:"promo_name"`
	Duration         int                `json:"duration"`
	PromoStartDate   string             `json:"promo_start_date"`
	PromoEndDate     string             `json:"promo_end_date"`
	IsActive         bool               `json:"is_active"`
	PromoType        string             `json:"promo_type"`
	PromoDetail      entity.PromoDetail `json:"promo_detail"`
	PromoDescription string             `json:"promo_description"`
}
