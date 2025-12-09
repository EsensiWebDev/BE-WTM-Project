package notification_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (nu *NotificationUsecase) ReadNotification(ctx context.Context, req *notifdto.ReadNotificationRequest) error {
	userCtx, err := nu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return fmt.Errorf("user context is nil")
	}

	agentID := userCtx.ID

	filterReq := filter.NotifFilter{
		UserID: agentID,
	}

	notifs, _, err := nu.notifRepo.GetNotificationsByUserID(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "get notification fail", err.Error())
		return err
	}

	if len(notifs) == 0 {
		logger.Error(ctx, "no notification found")
		return nil
	}

	var selectedIDs []uint
	if req.Type != "all" {
		for _, notif := range notifs {
			if req.ID == notif.ID {
				selectedIDs = append(selectedIDs, notif.ID)
				break
			}
		}
	} else {
		for _, notif := range notifs {
			selectedIDs = append(selectedIDs, notif.ID)
		}
	}

	if err := nu.notifRepo.ReadNotification(ctx, selectedIDs); err != nil {
		logger.Error(ctx, "read notification fail", err.Error())
		return err
	}

	return nil
}
