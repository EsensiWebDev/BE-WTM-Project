package promo_group_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) GetPromoGroups(ctx context.Context, search string, limit, page int) ([]entity.PromoGroup, int64, error) {
	var promoGroups []model.PromoGroup
	var total int64

	query := pgr.db.DB.WithContext(ctx).
		Select("id, name").
		Model(&model.PromoGroup{})

	if strings.TrimSpace(search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(search)
		query = query.Where("name ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		logger.Error(ctx, "Error counting promo groups", err.Error())
		return nil, total, err
	}

	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	err = query.Find(&promoGroups).Error
	if err != nil {
		logger.Error(ctx, "Error finding promo groups", err.Error())
		return nil, total, err
	}

	var promoGroupsEntity []entity.PromoGroup
	if err := utils.CopyPatch(&promoGroupsEntity, &promoGroups); err != nil {
		logger.Error(ctx, "Error copying promo groups model to entity", err.Error())
		return nil, total, err
	}

	return promoGroupsEntity, total, nil
}
