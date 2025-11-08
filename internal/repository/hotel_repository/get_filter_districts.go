package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterDistricts(ctx context.Context, filter filter.HotelFilterForAgent) ([]string, error) {
	db := hr.db.GetTx(ctx)

	// 1️⃣ Subquery: hotel_id + min_price filtered
	minPriceSub := db.
		Select("h.id AS hotel_id, MIN(rp.price) AS min_price").
		Table("hotels AS h").
		Joins("JOIN room_types rt ON h.id = rt.hotel_id").
		Joins("JOIN room_prices rp ON rt.id = rp.room_type_id").
		Group("h.id")

	if filter.PriceMin != nil && filter.PriceMax != nil {
		minPriceSub = minPriceSub.
			Having("MIN(rp.price) BETWEEN ? AND ?", *filter.PriceMin, *filter.PriceMax)
	} else if filter.PriceMin != nil {
		minPriceSub = minPriceSub.
			Having("MIN(rp.price) > ?", *filter.PriceMin)
	} else if filter.PriceMax != nil {
		minPriceSub = minPriceSub.
			Having("MIN(rp.price) < ?", *filter.PriceMax)
	}

	// 2️⃣ Final Query with EXISTS filter
	query := db.Table("hotels AS h").
		Select(`h.addr_sub_district`).
		Joins("JOIN (?) AS mp ON mp.hotel_id = h.id", minPriceSub)

	if len(filter.Ratings) > 0 {
		query = query.Where("h.rating IN ?", filter.Ratings)
	}

	if len(filter.TotalBedrooms) > 0 || len(filter.BedTypeIDs) > 0 {
		existsClause := db.
			Table("room_types rt").
			Select("1").
			Joins("JOIN bed_type_rooms btr ON btr.room_type_id = rt.id").
			Joins("JOIN bed_types bt ON bt.id = btr.bed_type_id").
			Where("rt.hotel_id = h.id") // <- penting, link ke outer "hotels AS h"

		if len(filter.TotalBedrooms) > 0 {
			existsClause = existsClause.Where("rt.total_unit IN ?", filter.TotalBedrooms)
		}
		if len(filter.BedTypeIDs) > 0 {
			existsClause = existsClause.Where("bt.name IN ?", filter.BedTypeIDs)
		}

		query = query.Where("EXISTS (?)", existsClause)
	}

	// Search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(h.name) ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	// Select distinct sub-districts
	var districts []string
	if err := query.Distinct().Pluck("addr_sub_district", &districts).Error; err != nil {
		logger.Error(ctx, "Error fetching filter districts", err.Error())
		return nil, err
	}

	// If no districts found, return an empty slice
	if len(districts) == 0 {
		return []string{}, nil
	}

	return districts, nil
}
