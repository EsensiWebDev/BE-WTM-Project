package notification_usecase

import "wtm-backend/internal/domain"

type NotificationUsecase struct {
	notifRepo  domain.NotificationRepository
	middleware domain.Middleware
	dbTrx      domain.DatabaseTransaction
}

func NewNotificationUsecase(notifRepo domain.NotificationRepository, middleware domain.Middleware, dbTrx domain.DatabaseTransaction) *NotificationUsecase {
	return &NotificationUsecase{
		notifRepo:  notifRepo,
		middleware: middleware,
		dbTrx:      dbTrx,
	}
}
