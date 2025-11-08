package promogroupdto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"wtm-backend/internal/dto"
)

type ListPromoGroupPromosRequest struct {
	ID                    uint `json:"id" form:"id"`
	dto.PaginationRequest `json:",inline"`
}

type ListPromoGroupPromosResponse struct {
	Promos []ListPromoGroupPromosData `json:"promos"`
}

type ListPromoGroupPromosData struct {
	PromoID        uint   `json:"promo_id"`
	PromoName      string `json:"promo_name"`
	PromoCode      string `json:"promo_code"`
	PromoStartDate string `json:"promo_start_date"`
	PromoEndDate   string `json:"promo_end_date"`
}

func (r *ListPromoGroupPromosRequest) Validate() error {
	return validation.ValidateStruct(r, validation.Field(&r.ID, validation.Required.Error("Promo group Id is required")))
}
