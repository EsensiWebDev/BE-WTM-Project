package promo_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) GetPromosWithHotels(ctx context.Context, filterReq *filter.PromoFilter) ([]entity.Promo, int64, error) {
	db := pr.db.GetTx(ctx)

	var promos []model.Promo
	var total int64

	// Base query sesuai SQL yang kamu tulis
	query := db.WithContext(ctx).
		Model(&model.Promo{}).
		Preload("PromoRoomTypes").
		Preload("PromoRoomTypes.RoomType").
		Preload("PromoRoomTypes.RoomType.Hotel").
		Where("is_active = ?", true).
		Where("id IN (?)",
			db.Table("detail_promo_groups").
				Select("promo_id").
				Where("promo_group_id IN (?)",
					db.Table("users").
						Select("promo_group_id").
						Where("id = ?", filterReq.AgentID),
				),
		)

	// Search filter
	if filterReq.Search != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filterReq.Search)
		query = query.Where("promos.name ILIKE ? ", "%"+safeSearch+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting promos", err.Error())
		return nil, total, err
	}

	// Pagination
	if filterReq.Limit > 0 {
		if filterReq.Page < 1 {
			filterReq.Page = 1
		}
		offset := (filterReq.Page - 1) * filterReq.Limit
		query = query.Limit(filterReq.Limit).Offset(offset)
	}

	// Fetch promos
	if err := query.Find(&promos).Error; err != nil {
		logger.Error(ctx, "Error finding promos", err.Error())
		return nil, total, err
	}

	// Convert to entity
	var promoEntities []entity.Promo
	if err := utils.CopyStrict(&promoEntities, &promos); err != nil {
		logger.Error(ctx, "Error copying promos model to entity", err.Error())
		return nil, total, err
	}

	// Map detail + promo type
	for i, promo := range promos {
		for i2, roomType := range promo.PromoRoomTypes {
			promoEntities[i].PromoRoomTypes[i2].RoomTypeName = roomType.RoomType.Name
			promoEntities[i].PromoRoomTypes[i2].HotelName = roomType.RoomType.Hotel.Name
			promoEntities[i].PromoRoomTypes[i2].Province = roomType.RoomType.Hotel.AddrProvince
			promoEntities[i].PromoRoomTypes[i2].TotalNights = roomType.TotalNights
		}

		// Set Duration from PromoRoomType.TotalNights if available
		// Use the first PromoRoomType's TotalNights as the promo's Duration
		// (assuming all room types for a promo have the same TotalNights requirement)
		if len(promo.PromoRoomTypes) > 0 && promo.PromoRoomTypes[0].TotalNights > 0 {
			promoEntities[i].Duration = promo.PromoRoomTypes[0].TotalNights
		}
	}

	return promoEntities, total, nil
}
