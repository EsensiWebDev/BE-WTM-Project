package currency_repository

import (
	"wtm-backend/internal/infrastructure/database"
)

type CurrencyRepository struct {
	db *database.DBPostgre
}

func NewCurrencyRepository(db *database.DBPostgre) *CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}
