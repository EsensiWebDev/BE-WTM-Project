package hotel_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterBedTypes(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterBedTypeHotel, error) {
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
		Select(`bt.id AS bed_type_id, bt.name AS bed_type, COUNT(DISTINCT h.id) AS count`).
		Joins("JOIN (?) AS mp ON mp.hotel_id = h.id", minPriceSub).
		Joins("JOIN room_types rt ON h.id = rt.hotel_id").
		Joins("JOIN bed_type_rooms btr ON btr.room_type_id = rt.id").
		Joins("JOIN bed_types bt ON btr.bed_type_id = bt.id")

	if len(filter.TotalBedrooms) > 0 {
		query = query.Where("rt.total_unit IN ?", filter.TotalBedrooms)
	}

	if len(filter.Cities) > 0 {
		query = query.Where("h.addr_city IN ?", filter.Cities)
	}
	if len(filter.Ratings) > 0 {
		query = query.Where("h.rating IN ?", filter.Ratings)
	}

	// Search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(h.name) ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	// Group by bed type
	query = query.Group("bt.id, bt.name")

	// Execute the query
	var results []entity.FilterBedTypeHotel
	if err := query.Scan(&results).Error; err != nil {
		logger.Error(ctx, "Error fetching filter bed types", err.Error())
		return nil, fmt.Errorf("error fetching filter bed types: %s", err.Error())
	}

	if len(results) == 0 {
		logger.Info(ctx, "No bed types found for the given filters")
		return nil, nil
	}

	return results, nil

}
