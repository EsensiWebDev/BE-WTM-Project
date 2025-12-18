package currency_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (cr *CurrencyRepository) CreateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error) {
	db := cr.db.GetTx(ctx)

	var currencyModel model.Currency
	if err := utils.CopyStrict(&currencyModel, currency); err != nil {
		logger.Error(ctx, "Error copying currency entity to model", err.Error())
		return nil, err
	}

	if err := db.WithContext(ctx).Create(&currencyModel).Error; err != nil {
		if cr.db.ErrDuplicateKey(ctx, err) {
			logger.Warn(ctx, "Currency already exists with code", currency.Code)
			return nil, err
		}
		logger.Error(ctx, "Error creating currency", err.Error())
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
