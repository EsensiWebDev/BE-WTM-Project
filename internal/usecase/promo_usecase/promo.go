package promo_usecase

import (
	"wtm-backend/internal/domain"
)

type PromoUsecase struct {
	promoRepo  domain.PromoRepository
	dbTrx      domain.DatabaseTransaction
	middleware domain.Middleware
}

func NewPromoUsecase(promoRepo domain.PromoRepository, dbTrx domain.DatabaseTransaction, middleware domain.Middleware) *PromoUsecase {
	return &PromoUsecase{
		promoRepo:  promoRepo,
		dbTrx:      dbTrx,
		middleware: middleware,
	}
}
