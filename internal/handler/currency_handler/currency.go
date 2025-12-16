package currency_handler

import (
	"wtm-backend/internal/domain"
)

type CurrencyHandler struct {
	currencyUsecase domain.CurrencyUsecase
}

func NewCurrencyHandler(currencyUsecase domain.CurrencyUsecase) *CurrencyHandler {
	return &CurrencyHandler{
		currencyUsecase: currencyUsecase,
	}
}
