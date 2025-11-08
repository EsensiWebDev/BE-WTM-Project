package notification_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (nu *NotificationUsecase) ListNotifications(ctx context.Context, req *notifdto.ListNotificationsRequest) (*notifdto.ListNotificationsResponse, error) {

	userCtx, err := nu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", "error", err)
		return nil, fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return nil, fmt.Errorf("user context is nil")
	}

	userID := userCtx.ID

	filterReq := filter.NotifFilter{
		UserID: userID,
	}
	filterReq.PaginationRequest = req.PaginationRequest

	notifications, total, err := nu.notifRepo.GetNotificationsByUserID(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get notifications by user Id", "error", err)
		return nil, fmt.Errorf("failed to get notifications by user Id: %s", err.Error())
	}
	response := &notifdto.ListNotificationsResponse{
		Notifications: notifications,
		Total:         total,
	}

	return response, nil
}
