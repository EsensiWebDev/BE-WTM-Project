package promo_group_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) GetUnassignedPromos(ctx context.Context, filterReq *filter.PromoGroupFilter) ([]entity.Promo, int64, error) {
	db := pgr.db.GetTx(ctx)

	var promos []model.Promo
	var total int64
	query := db.WithContext(ctx).Model(&model.Promo{}).Select("id, name")

	// 🔍 Filter: belum berelasi dengan promoGroupID
	query = query.Where("NOT EXISTS (?)",
		db.Table("detail_promo_groups").
			Select("1").
			Where("detail_promo_groups.promo_id = promos.id").
			Where("detail_promo_groups.promo_group_id = ?", filterReq.PromoGroupID),
	)

	// 🔍 Filter: search by name
	if strings.TrimSpace(filterReq.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filterReq.Search)
		query = query.Where("LOWER(promos.name) ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	// 🔢 Count total
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting unassigned promos:", err.Error())
		return nil, 0, err
	}

	// 📦 Pagination
	limit := filterReq.Limit
	page := filterReq.Page
	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	// 🏁 Execute query
	if err := query.Find(&promos).Error; err != nil {
		logger.Error(ctx, "Error fetching unassigned promos:", err.Error())
		return nil, 0, err
	}

	// 🔄 Convert to entity
	var promoEntities []entity.Promo
	if err := utils.CopyStrict(&promoEntities, promos); err != nil {
		logger.Error(ctx, "Error converting promo models to entities:", err.Error())
		return nil, 0, err
	}

	return promoEntities, total, nil
}
