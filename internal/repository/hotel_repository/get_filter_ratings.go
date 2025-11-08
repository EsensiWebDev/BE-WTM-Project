package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterRatings(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterRatingHotel, error) {
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
		Select(`h.rating, count(DISTINCT h.id) AS count`).
		Joins("JOIN (?) AS mp ON mp.hotel_id = h.id", minPriceSub)

	if len(filter.Cities) > 0 {
		query = query.Where("h.addr_city IN ?", filter.Cities)
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

	// Execute
	var ratings []entity.FilterRatingHotel
	if err := query.Group("h.rating").Order("h.rating").Find(&ratings).Error; err != nil {
		logger.Error(ctx, "Error fetching filter ratings", err.Error())
		return nil, err
	}

	// Check if no ratings found
	if len(ratings) == 0 {
		logger.Info(ctx, "No ratings found for the given filter criteria")
		return nil, nil
	}

	return ratings, nil
}
