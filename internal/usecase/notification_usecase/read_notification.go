package notification_usecase

import (
	"context"
	"wtm-backend/pkg/logger"
)

func (nu *NotificationUsecase) ReadNotification(ctx context.Context, id int64) error {
	if err := nu.notifRepo.ReadNotification(ctx, id); err != nil {
		logger.Error(ctx, "read notification fail", err.Error())
		return err
	}

	return nil
}
