package email_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (er *EmailRepository) GetEmailLogByID(ctx context.Context, id uint) (*entity.EmailLog, error) {
	db := er.db.GetTx(ctx)

	var emailLog model.EmailLog
	if err := db.WithContext(ctx).
		Preload("EmailTemplate").
		Where("id = ?", id).
		First(&emailLog).Error; err != nil {
		logger.Error(ctx, "failed to get email log by id", err.Error())
		return nil, err
	}

	var result entity.EmailLog
	if err := utils.CopyStrict(&result, &emailLog); err != nil {
		logger.Error(ctx, "failed to copy email log model to entity", err.Error())
		return nil, err
	}

	if emailLog.Meta != nil {
		var meta entity.MetadataEmailLog
		if err := json.Unmarshal(emailLog.Meta, &meta); err != nil {
			logger.Error(ctx, "failed to unmarshal email log meta", err.Error())
			return nil, err
		}
		result.Meta = &meta
	}

	// Safely get email type from template
	result.EmailType = getEmailTypeFromTemplate(ctx, emailLog)

	return &result, nil
}
