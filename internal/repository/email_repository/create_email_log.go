package email_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (er *EmailRepository) CreateEmailLog(ctx context.Context, log *entity.EmailLog) error {
	db := er.db.GetTx(ctx)

	var modelEmailLog model.EmailLog
	if err := utils.CopyStrict(&modelEmailLog, log); err != nil {
		logger.Error(ctx, "failed to copy email log entity to model", err.Error())
		return err
	}

	metaJSON, err := json.Marshal(log.Meta)
	if err != nil {
		logger.Error(ctx, "failed to marshal metadata to JSON", err.Error())
		return err
	}
	modelEmailLog.Meta = metaJSON

	modelEmailLog.StatusID = constant.StatusEmailPendingID

	if err := db.WithContext(ctx).Create(&modelEmailLog).Error; err != nil {
		logger.Error(ctx, "failed to create email log", err.Error())
		return err
	}

	return nil
}
