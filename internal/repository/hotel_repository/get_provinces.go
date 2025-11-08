package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetProvinces(ctx context.Context, filter *filter.DefaultFilter) ([]string, int64, error) {
	db := hr.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Hotel{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(addr_province) ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
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

	var provinces []string
	if err := query.
		Select("DISTINCT LOWER(addr_province) AS addr_province").
		Pluck("addr_province", &provinces).Error; err != nil {
		logger.Error(ctx, "Error fetching facilities", err.Error())
		return nil, total, err
	}

	if len(provinces) == 0 {
		logger.Info(ctx, "No provinces found")
		return nil, total, nil
	}

	return provinces, total, nil
}
