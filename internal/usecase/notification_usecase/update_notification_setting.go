package notification_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/notifdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (nu *NotificationUsecase) UpdateNotificationSetting(ctx context.Context, req *notifdto.UpdateNotificationSettingRequest) error {
	user, err := nu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error generating user from context", err.Error())
		return err
	}

	if user == nil {
		logger.Error(ctx, "User not found in context")
		return err
	}

	return nu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := nu.notifRepo.DisableNotificationSettingsByChannel(txCtx, user.ID, req.Channel); err != nil {
			logger.Error(ctx, "Error disabling notification settings by channel", err.Error())
			return err
		}

		if !req.IsEnable {
			logger.Info(ctx, "Disabled notification settings by channel")
			return nil
		}

		var typesToEnable []string
		switch req.Type {
		case constant.ConstBooking, constant.ConstReject:
			typesToEnable = []string{req.Type}
		case constant.ConstAll:
			typesToEnable = []string{constant.ConstBooking, constant.ConstReject}
		default:
			return fmt.Errorf("invalid type: %s", req.Type)
		}

		if err := nu.notifRepo.EnableNotificationSettings(txCtx, user.ID, req.Channel, typesToEnable); err != nil {
			logger.Error(ctx, "Error enabling notification settings", err.Error())
			return err
		}

		return nil
	})
}
