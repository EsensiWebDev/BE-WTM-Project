package hotel_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetFilterBedTypes(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterBedTypeHotel, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var hotelConditions []string
	var roomConditions []string
	var priceHaving string

	// ⚠️ TIDAK include BedTypeIDs (karena ini fungsi untuk get available bed types)

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

	// Additional JOINs untuk bed types
	additionalJoins := `
		JOIN room_types rt ON h.id = rt.hotel_id
		JOIN bed_type_rooms btr ON btr.room_type_id = rt.id
		JOIN bed_types bt ON btr.bed_type_id = bt.id
	`

	// Build query
	query := hr.buildBaseHotelQuery(
		`SELECT bt.id AS bed_type_id, bt.name AS bed_type, COUNT(DISTINCT h.id) AS count`,
		roomConditions,
		priceHaving,
		hotelConditions,
		additionalJoins,
		"GROUP BY bt.id, bt.name",
		"ORDER BY bt.name ASC",
	)

	// 🔍 Execute query
	var results []entity.FilterBedTypeHotel
	if err := db.Raw(query, args...).Debug().Scan(&results).Error; err != nil {
		logger.Error(ctx, "Error fetching filter bed types (raw)", err.Error())
		return nil, fmt.Errorf("error fetching filter bed types: %s", err.Error())
	}

	if len(results) == 0 {
		logger.Info(ctx, "No bed types found for the given filters")
		return []entity.FilterBedTypeHotel{}, nil
	}

	return results, nil
}
