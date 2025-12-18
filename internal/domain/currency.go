package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
)

type CurrencyRepository interface {
	GetAllCurrencies(ctx context.Context) ([]entity.Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*entity.Currency, error)
	GetCurrencyByID(ctx context.Context, id uint) (*entity.Currency, error)
	CreateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error)
	UpdateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error)
	GetActiveCurrencies(ctx context.Context) ([]entity.Currency, error)
}

type CurrencyUsecase interface {
	GetAllCurrencies(ctx context.Context) ([]entity.Currency, error)
	GetActiveCurrencies(ctx context.Context) ([]entity.Currency, error)
	CreateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error)
	UpdateCurrency(ctx context.Context, currency *entity.Currency) (*entity.Currency, error)
}
