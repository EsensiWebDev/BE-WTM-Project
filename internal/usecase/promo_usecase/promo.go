package promo_usecase

import (
	"wtm-backend/internal/domain"
)

type PromoUsecase struct {
	promoRepo domain.PromoRepository
	dbTrx     domain.DatabaseTransaction
}

func NewPromoUsecase(promoRepo domain.PromoRepository, dbTrx domain.DatabaseTransaction) *PromoUsecase {
	return &PromoUsecase{
		promoRepo: promoRepo,
		dbTrx:     dbTrx,
	}
}
