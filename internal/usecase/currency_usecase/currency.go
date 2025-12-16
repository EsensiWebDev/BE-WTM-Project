package currency_usecase

import (
	"wtm-backend/internal/domain"
)

type CurrencyUsecase struct {
	currencyRepo domain.CurrencyRepository
}

func NewCurrencyUsecase(currencyRepo domain.CurrencyRepository) *CurrencyUsecase {
	return &CurrencyUsecase{
		currencyRepo: currencyRepo,
	}
}
