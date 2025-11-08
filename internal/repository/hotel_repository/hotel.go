package hotel_repository

import (
	"wtm-backend/internal/infrastructure/database"
)

type HotelRepository struct {
	db *database.DBPostgre
}

func NewHotelRepository(db *database.DBPostgre) *HotelRepository {
	return &HotelRepository{
		db: db,
	}
}
