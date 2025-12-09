package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterDistricts(ctx context.Context, filter filter.HotelFilterForAgent) ([]string, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var hotelConditions []string
	var roomConditions []string
	var priceHaving string

	// 🔍 Filter bed type
	if len(filter.BedTypeIDs) > 0 {
		roomConditions = append(roomConditions, "bt.id IN ?")
		args = append(args, filter.BedTypeIDs)
	}

	// 🔍 Filter total bedrooms
	if len(filter.TotalBedrooms) > 0 {
		roomConditions = append(roomConditions, "rt.total_unit IN ?")
		args = append(args, filter.TotalBedrooms)
	}

	// 🔍 Filter min guest
	if filter.MinGuest > 0 {
		roomConditions = append(roomConditions, "rt.max_occupancy >= ?")
		args = append(args, filter.MinGuest)
	}

	// 🔍 Filter promo
	if filter.PromoID > 0 {
		roomConditions = append(roomConditions, `
			EXISTS (
				SELECT 1 
				FROM promo_room_types prt
				JOIN promos p ON prt.promo_id = p.id
				WHERE prt.room_type_id = rt.id
				AND prt.promo_id = ?
				AND p.is_active = true
			)
		`)
		args = append(args, filter.PromoID)
	}

	// 🔍 Filter availability
	if filter.DateFrom != nil && filter.DateTo != nil {
		roomConditions = append(roomConditions,
			"NOT EXISTS (SELECT 1 FROM room_unavailables ru WHERE ru.room_type_id = rt.id AND ru.date BETWEEN ? AND ?)")
		args = append(args, *filter.DateFrom, *filter.DateTo)
	}

	// 🔍 Filter harga
	if filter.PriceMin != nil && filter.PriceMax != nil {
		priceHaving = "HAVING MIN(rp.price) BETWEEN ? AND ?"
		args = append(args, *filter.PriceMin, *filter.PriceMax)
	} else if filter.PriceMin != nil {
		priceHaving = "HAVING MIN(rp.price) >= ?"
		args = append(args, *filter.PriceMin)
	} else if filter.PriceMax != nil {
		priceHaving = "HAVING MIN(rp.price) <= ?"
		args = append(args, *filter.PriceMax)
	}

	// 🔍 Filter province
	if filter.Province != nil && strings.TrimSpace(*filter.Province) != "" {
		hotelConditions = append(hotelConditions, "h.addr_province = ?")
		args = append(args, *filter.Province)
	}

	// 🔍 Filter rating
	if len(filter.Ratings) > 0 {
		hotelConditions = append(hotelConditions, "h.rating IN ?")
		args = append(args, filter.Ratings)
	}

	// 🔍 Filter nama hotel (search)
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		hotelConditions = append(hotelConditions, "LOWER(h.name) ILIKE ?")
		args = append(args, "%"+safeSearch+"%")
	}

	// 🔍 Filter status hotel
	hotelConditions = append(hotelConditions, "h.status_id = ?")
	args = append(args, constant.StatusHotelApprovedID)

	// ⚠️ TIDAK include Cities (karena ini fungsi untuk get list cities)

	// Build query
	query := hr.buildBaseHotelQuery(
		`SELECT DISTINCT h.addr_city`,
		roomConditions,
		priceHaving,
		hotelConditions,
		"", // no additional joins
		"", // no group by
		"ORDER BY h.addr_city ASC",
	)

	// 🔍 Execute query
	var districts []string
	if err := db.Raw(query, args...).Scan(&districts).Error; err != nil {
		logger.Error(ctx, "Error fetching filter districts (raw)", err.Error())
		return nil, err
	}

	if len(districts) == 0 {
		return []string{}, nil
	}

	return districts, nil
}
