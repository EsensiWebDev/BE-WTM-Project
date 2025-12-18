package currency_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (cu *CurrencyUsecase) UpdateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error) {
	// Check if currency exists
	existing, err := cu.currencyRepo.GetCurrencyByID(ctx, currency.ID)
	if err != nil {
		logger.Error(ctx, "Error checking existing currency", err.Error())
		return nil, err
	}
	if existing == nil {
		logger.Warn(ctx, "Currency not found", currency.ID)
		return nil, fmt.Errorf("currency with id %d not found", currency.ID)
	}

	// Update currency (code cannot be changed)
	currency.Code = existing.Code
	updated, err := cu.currencyRepo.UpdateCurrency(ctx, currency)
	if err != nil {
		logger.Error(ctx, "Error updating currency", err.Error())
		return nil, err
	}

	return updated, nil
}
