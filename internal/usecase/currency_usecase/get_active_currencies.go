package currency_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (cu *CurrencyUsecase) GetActiveCurrencies(ctx context.Context) ([]entity.Currency, error) {
	currencies, err := cu.currencyRepo.GetActiveCurrencies(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting active currencies", err.Error())
		return nil, err
	}
	return currencies, nil
}
