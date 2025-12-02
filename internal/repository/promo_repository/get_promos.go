package promo_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) GetPromos(ctx context.Context, filterReq *filter.DefaultFilter) ([]entity.Promo, int64, error) {
	db := pr.db.GetTx(ctx)

	var promos []model.Promo
	var total int64

	query := db.WithContext(ctx).
		Model(&model.Promo{}).
		Preload("PromoType").
		Preload("PromoRoomTypes")

	if filterReq.Search != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filterReq.Search)
		query = query.Where("promos.name ILIKE ? ", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting promos", err.Error())
		return nil, total, err
	}

	if filterReq.Limit > 0 {
		if filterReq.Page < 1 {
			filterReq.Page = 1
		}
		offset := (filterReq.Page - 1) * filterReq.Limit
		query = query.Limit(filterReq.Limit).Offset(offset)
	}

	query = query.Order("promos.created_at DESC")

	if err := query.Find(&promos).Error; err != nil {
		logger.Error(ctx, "Error finding promos", err.Error())
		return nil, total, err
	}

	var promoEntities []entity.Promo
	if err := utils.CopyStrict(&promoEntities, &promos); err != nil {
		logger.Error(ctx, "Error copying promos model to entity", err.Error())
		return nil, total, err
	}

	for i, promo := range promos {
		promoEntities[i].PromoTypeName = promo.PromoType.Name
		var detailPromo entity.PromoDetail
		if err := json.Unmarshal(promo.Detail, &detailPromo); err != nil {
			logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
		}
		promoEntities[i].Detail = detailPromo
		promoEntities[i].ExternalID = promo.ExternalID.ExternalID
	}

	return promoEntities, total, nil
}
