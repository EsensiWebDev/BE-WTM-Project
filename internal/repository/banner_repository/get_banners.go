package banner_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BannerRepository) GetBanners(ctx context.Context, filter *filter.BannerFilter) ([]entity.Banner, int64, error) {
	db := br.db.GetTx(ctx)

	query := db.WithContext(ctx).Model(&model.Banner{})

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(title) ILIKE ? ", "%"+safeSearch+"%")
	}

	// Apply filter for active banners
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
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
		query = query.Offset(offset).Limit(filter.Limit)
	}

	//Apply sorting
	query = query.Order("display_order DESC")

	// Execute the query to fetch active banners
	var banners []model.Banner
	if err := query.Find(&banners).Error; err != nil {
		logger.Error(ctx, "Error fetching active banners", err.Error())
		return nil, total, err
	}

	// Mapping
	var activeBanners []entity.Banner
	if err := utils.CopyStrict(&activeBanners, &banners); err != nil {
		logger.Error(ctx, "Error copying banners to entity", err.Error())
		return nil, total, err
	}
	for i, banner := range banners {
		activeBanners[i].ExternalID = banner.ExternalID.ExternalID
	}

	return activeBanners, total, nil
}
