package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetBedTypes(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error) {
	db := hr.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.BedType{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("name ILIKE ? ", "%"+safeSearch+"%")
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting bed types", err.Error())
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

	var bedTypes []string
	if err := query.Pluck("name", &bedTypes).Error; err != nil {
		logger.Error(ctx, "Error fetching bed types", err.Error())
		return nil, total, err
	}

	return bedTypes, total, nil
}
