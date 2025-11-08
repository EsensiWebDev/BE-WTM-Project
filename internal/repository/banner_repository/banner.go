package banner_repository

import (
	"wtm-backend/internal/infrastructure/database"
)

type BannerRepository struct {
	db *database.DBPostgre
}

func NewBannerRepository(db *database.DBPostgre) *BannerRepository {
	return &BannerRepository{
		db: db,
	}
}
