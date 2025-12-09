package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterPricing(ctx context.Context, filter filter.HotelFilterForAgent) (*entity.FilterRangePrice, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var hotelConditions []string
	var roomConditions []string

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

	// ⚠️ TIDAK include PriceMin/PriceMax (karena ini fungsi untuk get range)

	// 🔍 Filter province
	if filter.Province != nil && strings.TrimSpace(*filter.Province) != "" {
		hotelConditions = append(hotelConditions, "h.addr_province = ?")
		args = append(args, *filter.Province)
	}

	// 🔍 Filter kota
	if len(filter.Cities) > 0 {
		hotelConditions = append(hotelConditions, "h.addr_city IN ?")
		args = append(args, filter.Cities)
	}

	// 🔍 Filter rating
	if len(filter.Ratings) > 0 {
		hotelConditions = append(hotelConditions, "h.rating IN ?")
		args = append(args, filter.Ratings)
	}

	// 🔍 Filter search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		hotelConditions = append(hotelConditions, "LOWER(h.name) ILIKE ?")
		args = append(args, "%"+safeSearch+"%")
	}

	// 🔍 Filter status hotel
	hotelConditions = append(hotelConditions, "h.status_id = ?")
	args = append(args, constant.StatusHotelApprovedID)

	// Build query
	query := hr.buildBaseHotelQuery(
		`SELECT MIN(mp.min_price) AS min_price, MAX(mp.min_price) AS max_price`,
		roomConditions,
		"", // no price HAVING
		hotelConditions,
		"", // no additional joins
		"", // no group by
		"", // no order by
	)

	// 🔍 Execute query
	var result entity.FilterRangePrice
	if err := db.Raw(query, args...).Scan(&result).Error; err != nil {
		logger.Error(ctx, "Error fetching range price (raw)", err.Error())
		return nil, err
	}

	return &result, nil
}
