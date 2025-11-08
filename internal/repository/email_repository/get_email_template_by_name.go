package email_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (er *EmailRepository) GetEmailTemplateByName(ctx context.Context, name string) (*entity.EmailTemplate, error) {
	db := er.db.GetTx(ctx)

	var emailTemplate model.EmailTemplate
	if err := db.WithContext(ctx).Where("name = ?", name).First(&emailTemplate).Error; err != nil {
		logger.Error(ctx, "Failed to get email template by name", err.Error())
		return nil, err
	}

	var result entity.EmailTemplate
	if err := utils.CopyStrict(&result, &emailTemplate); err != nil {
		logger.Error(ctx, "Failed to copy email template model to entity", err.Error())
		return nil, err
	}

	return &result, nil
}
