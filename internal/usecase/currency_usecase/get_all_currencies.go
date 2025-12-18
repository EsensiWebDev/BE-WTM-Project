package currency_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (cu *CurrencyUsecase) GetAllCurrencies(ctx context.Context) ([]entity.Currency, error) {
	currencies, err := cu.currencyRepo.GetAllCurrencies(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting all currencies", err.Error())
		return nil, err
	}
	return currencies, nil
}
