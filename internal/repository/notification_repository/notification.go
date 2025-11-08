package notification_repository

import "wtm-backend/internal/infrastructure/database"

type NotificationRepository struct {
	db *database.DBPostgre
}

func NewNotificationRepository(db *database.DBPostgre) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}
