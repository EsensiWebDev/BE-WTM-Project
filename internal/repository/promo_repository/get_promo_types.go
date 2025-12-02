package promo_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) GetPromoTypes(ctx context.Context, filter *filter.DefaultFilter) ([]entity.PromoType, int64, error) {
	db := pr.db.GetTx(ctx)

	var promoTypes []model.PromoType
	var total int64

	query := db.WithContext(ctx).
		Model(&model.PromoType{}).
		Select("id, name")

	if filter.Search != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("name ILIKE ? ", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting promo types", err.Error())
		return nil, total, err
	}

	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	if err := query.Find(&promoTypes).Error; err != nil {
		logger.Error(ctx, "Error finding promo types", err.Error())
		return nil, total, err
	}

	var promoTypeEntity []entity.PromoType
	if err := utils.CopyPatch(&promoTypeEntity, &promoTypes); err != nil {
		logger.Error(ctx, "Error copying promo types model to entity", err.Error())
		return nil, total, err
	}

	return promoTypeEntity, total, nil
}
