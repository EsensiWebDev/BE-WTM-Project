package user_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetStatusUsers(ctx context.Context, filter *filter.DefaultFilter) ([]entity.StatusUser, int64, error) {
	db := ur.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.StatusUser{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(status) ILIKE ? ", "%"+safeSearch+"%")
	}

	// Apply filter for active banners
	query = query.Where("is_active = ?", true)

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting facilities", err.Error())
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

	// Fetch records
	var statusUsers []model.StatusUser
	if err := query.Find(&statusUsers).Error; err != nil {
		logger.Error(ctx, "Error fetching status users", err.Error())
		return nil, total, err
	}

	// Mapping
	var mappedStatusUsers []entity.StatusUser
	if err := utils.CopyPatch(&mappedStatusUsers, statusUsers); err != nil {
		logger.Error(ctx, "Error mapping status users", err.Error())
		return nil, total, err
	}

	return mappedStatusUsers, total, nil
}
