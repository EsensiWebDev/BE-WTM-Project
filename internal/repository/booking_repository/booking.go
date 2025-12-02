package booking_repository

import (
	"wtm-backend/internal/domain"
	"wtm-backend/internal/infrastructure/database"
)

type BookingRepository struct {
	db          *database.DBPostgre
	redisClient domain.RedisClient
}

func NewBookingRepository(db *database.DBPostgre, redisClient domain.RedisClient) *BookingRepository {
	return &BookingRepository{
		db:          db,
		redisClient: redisClient,
	}
}
