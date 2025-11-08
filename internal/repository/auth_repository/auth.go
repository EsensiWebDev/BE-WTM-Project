package auth_repository

import (
	"wtm-backend/internal/domain"
	"wtm-backend/internal/repository/driver"
)

type AuthRepository struct {
	redisClient domain.RedisClient
	db          driver.DBPostgre
}

func NewAuthRepository(redisClient domain.RedisClient, db driver.DBPostgre) *AuthRepository {
	return &AuthRepository{
		redisClient: redisClient,
		db:          db,
	}
}
