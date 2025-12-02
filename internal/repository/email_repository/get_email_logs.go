package email_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (er *EmailRepository) GetEmailLogs(ctx context.Context, filter filter.DefaultFilter) ([]entity.EmailLog, int64, error) {
	db := er.db.GetTx(ctx)
	query := db.WithContext(ctx).Model(&model.EmailLog{})

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "failed to count email logs", err.Error())
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// Apply sorting
	query = query.Order("created_at DESC")

	// Fetch results
	var emailLogs []model.EmailLog
	if err := query.Find(&emailLogs).Error; err != nil {
		logger.Error(ctx, "failed to get email logs", err.Error())
		return nil, 0, err
	}

	// Convert to entity.EmailLog slice
	var result []entity.EmailLog
	if err := utils.CopyStrict(&result, &emailLogs); err != nil {
		logger.Error(ctx, "failed to copy email logs", err.Error())
		return nil, 0, err
	}

	for i, emailLog := range emailLogs {
		var meta entity.MetadataEmailLog
		if err := json.Unmarshal(emailLog.Meta, &meta); err != nil {
			logger.Error(ctx, "failed to unmarshal email log meta", err.Error())
			return nil, 0, err
		}
		result[i].Meta = meta
	}

	return result, total, nil

}
