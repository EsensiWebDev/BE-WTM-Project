package currency_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (cr *CurrencyRepository) GetAllCurrencies(ctx context.Context) ([]entity.Currency, error) {
	db := cr.db.GetTx(ctx)

	var currencies []model.Currency
	if err := db.WithContext(ctx).Find(&currencies).Error; err != nil {
		logger.Error(ctx, "Error fetching all currencies", err.Error())
		return nil, err
	}

	var entityCurrencies []entity.Currency
	for _, currency := range currencies {
		var entityCurrency entity.Currency
		if err := utils.CopyStrict(&entityCurrency, currency); err != nil {
			logger.Error(ctx, "Error copying currency model to entity", err.Error())
			return nil, err
		}
		entityCurrency.ExternalID = currency.ExternalID.ExternalID
		entityCurrencies = append(entityCurrencies, entityCurrency)
	}

	return entityCurrencies, nil
}
