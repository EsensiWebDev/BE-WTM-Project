package user_repository

import (
	"wtm-backend/internal/infrastructure/database"
	"wtm-backend/internal/repository/driver"
)

type UserRepository struct {
	db driver.DBPostgre
}

func NewUserRepository(db *database.DBPostgre) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
