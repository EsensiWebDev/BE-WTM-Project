package promo_repository

import (
	"wtm-backend/internal/infrastructure/database"
)

type PromoRepository struct {
	db *database.DBPostgre
}

func NewPromoRepository(db *database.DBPostgre) *PromoRepository {
	return &PromoRepository{
		db: db,
	}
}
