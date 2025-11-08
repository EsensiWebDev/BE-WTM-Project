package email_repository

import "wtm-backend/internal/infrastructure/database"

type EmailRepository struct {
	db *database.DBPostgre
}

func NewEmailRepository(db *database.DBPostgre) *EmailRepository {
	return &EmailRepository{
		db: db,
	}
}
