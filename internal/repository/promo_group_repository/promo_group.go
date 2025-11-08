package promo_group_repository

import "wtm-backend/internal/infrastructure/database"

type PromoGroupRepository struct {
	db *database.DBPostgre
}

func NewPromoGroupRepository(db *database.DBPostgre) *PromoGroupRepository {
	return &PromoGroupRepository{
		db: db,
	}
}
