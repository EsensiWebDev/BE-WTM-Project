package promo_handler

import (
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/usecase/promo_usecase"
)

type PromoHandler struct {
	promoUsecase domain.PromoUsecase
	config       *config.Config
}

func NewPromoHandler(config *config.Config, promoUsecase *promo_usecase.PromoUsecase) *PromoHandler {
	return &PromoHandler{config: config, promoUsecase: promoUsecase}
}
