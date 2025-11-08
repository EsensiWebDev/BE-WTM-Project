package promo_group_handler

import "wtm-backend/internal/domain"

type PromoGroupHandler struct {
	promoGroupUsecase domain.PromoGroupUsecase
}

func NewPromoGroupHandler(promoGroupUsecase domain.PromoGroupUsecase) *PromoGroupHandler {
	return &PromoGroupHandler{
		promoGroupUsecase: promoGroupUsecase,
	}
}
