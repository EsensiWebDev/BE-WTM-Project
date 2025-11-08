package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFacilities(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error) {
	db := hr.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Facility{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("name ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

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

	var facilities []string
	if err := query.Pluck("name", &facilities).Error; err != nil {
		logger.Error(ctx, "Error fetching facilities", err.Error())
		return nil, total, err
	}

	return facilities, total, nil
}
