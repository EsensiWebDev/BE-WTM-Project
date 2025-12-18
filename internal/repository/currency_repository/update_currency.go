package currency_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (cr *CurrencyRepository) UpdateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error) {
	db := cr.db.GetTx(ctx)

	var currencyModel model.Currency
	if err := utils.CopyStrict(&currencyModel, currency); err != nil {
		logger.Error(ctx, "Error copying currency entity to model", err.Error())
		return nil, err
	}

	if err := db.WithContext(ctx).Model(&model.Currency{}).
		Where("id = ?", currency.ID).
		Updates(map[string]interface{}{
			"name":      currency.Name,
			"symbol":    currency.Symbol,
			"is_active": currency.IsActive,
		}).Error; err != nil {
		logger.Error(ctx, "Error updating currency", err.Error())
		return nil, err
	}

	// Fetch updated currency
	if err := db.WithContext(ctx).Where("id = ?", currency.ID).First(&currencyModel).Error; err != nil {
		logger.Error(ctx, "Error fetching updated currency", err.Error())
		return nil, err
	}

	var entityCurrency entity.Currency
	if err := utils.CopyStrict(&entityCurrency, currencyModel); err != nil {
		logger.Error(ctx, "Error copying currency model to entity", err.Error())
		return nil, err
	}
	entityCurrency.ExternalID = currencyModel.ExternalID.ExternalID

	return &entityCurrency, nil
}
