package promo_group_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) GetPromosByPromoGroupID(ctx context.Context, promoGroupID uint, search string, limit, page int) ([]entity.Promo, int64, error) {
	db := pgr.db.GetTx(ctx)

	var promos []model.Promo
	var total int64

	query := db.WithContext(ctx).
		Model(&model.Promo{}).
		Select("promos.id, promos.code, promos.name, promos.start_date, promos.end_date").
		Joins("JOIN detail_promo_groups ON detail_promo_groups.promo_id = promos.id").
		Where("detail_promo_groups.promo_group_id = ?", promoGroupID)

	if strings.TrimSpace(search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(search)
		query = query.Where("promos.name ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
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

	if err := query.Find(&promos).Error; err != nil {
		logger.Error(ctx, "Error finding promos by promo group Id", err.Error())
		return nil, total, err
	}

	entityPromo := make([]entity.Promo, 0, len(promos))
	if err := utils.CopyPatch(&entityPromo, &promos); err != nil {
		logger.Error(ctx, "Error copying promos model to entity", err.Error())
		return nil, total, err
	}

	return entityPromo, total, nil
}
