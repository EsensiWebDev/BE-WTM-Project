package currency_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/domain/entity"
	currencypkg "wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
)

func (cu *CurrencyUsecase) CreateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error) {
	// Normalize currency code
	currency.Code = currencypkg.NormalizeCurrencyCode(currency.Code)

	// Validate currency code
	if !currencypkg.ValidateCurrencyCode(currency.Code) {
		logger.Error(ctx, "Invalid currency code", currency.Code)
		return nil, fmt.Errorf("invalid currency code: %s", currency.Code)
	}

	// Check if currency already exists
	existing, err := cu.currencyRepo.GetCurrencyByCode(ctx, currency.Code)
	if err != nil {
		logger.Error(ctx, "Error checking existing currency", err.Error())
		return nil, err
	}
	if existing != nil {
		logger.Warn(ctx, "Currency already exists", currency.Code)
		return nil, fmt.Errorf("currency with code %s already exists", currency.Code)
	}

	// Create currency
	created, err := cu.currencyRepo.CreateCurrency(ctx, currency)
	if err != nil {
		logger.Error(ctx, "Error creating currency", err.Error())
		return nil, err
	}

	return created, nil
}
