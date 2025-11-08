package email_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (er *EmailRepository) UpdateEmailTemplate(ctx context.Context, template *entity.EmailTemplate) error {
	db := er.db.GetTx(ctx)

	var templateModel model.EmailTemplate
	if err := utils.CopyStrict(&templateModel, template); err != nil {
		logger.Error(ctx, "failed to copy email template entity to model", err.Error())
		return err
	}

	if err := db.WithContext(ctx).Save(templateModel).Error; err != nil {
		logger.Error(ctx, "failed to update email template", err.Error())
		return err
	}

	return nil

}
