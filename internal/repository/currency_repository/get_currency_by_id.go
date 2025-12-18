package currency_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (cr *CurrencyRepository) GetCurrencyByID(ctx context.Context, id uint) (*entity.Currency, error) {
	db := cr.db.GetTx(ctx)

	var currency model.Currency
	if err := db.WithContext(ctx).Where("id = ?", id).First(&currency).Error; err != nil {
		if cr.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "Currency not found with id", id)
			return nil, nil
		}
		logger.Error(ctx, "Error finding currency by id", err.Error())
		return nil, err
	}

	var entityCurrency entity.Currency
	if err := utils.CopyStrict(&entityCurrency, currency); err != nil {
		logger.Error(ctx, "Error copying currency model to entity", err.Error())
		return nil, err
	}
	entityCurrency.ExternalID = currency.ExternalID.ExternalID

	return &entityCurrency, nil
}
