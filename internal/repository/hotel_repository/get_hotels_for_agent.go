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

func (hr *HotelRepository) GetHotelsForAgent(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.CustomHotel, int64, error) {
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

	// 🔍 Filter nama hotel
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		hotelConditions = append(hotelConditions, "LOWER(h.name) ILIKE ? ")
		args = append(args, "%"+safeSearch+"%")
	}

	// 🔍 Filter status hotel
	hotelConditions = append(hotelConditions, "h.status_id = ?")
	args = append(args, constant.StatusHotelApprovedID)

	// Build base query (tanpa LastInternalID)
	baseQuery := hr.buildBaseHotelQuery(
		`SELECT h.id, h.name, h.addr_province, h.addr_city, h.addr_sub_district, h.photos, h.rating, h.created_at, mp.min_price`,
		roomConditions,
		priceHaving,
		hotelConditions,
		"", // no additional joins
		"", // no group by
		"", // order by ditambahkan nanti
	)

	// 🔢 Count total
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS sub", baseQuery)
	if err := db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		logger.Error(ctx, "Error counting hotels (raw)", err.Error())
		return nil, 0, err
	}

	// 📦 Build final query
	finalQuery := baseQuery

	// Tambahkan ORDER BY
	finalQuery += "\n\t\tORDER BY h.created_at DESC, h.id ASC"

	// Tambahkan LIMIT dan OFFSET
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		finalQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)

	}

	// 🔍 Execute main query
	var hotels []entity.CustomHotel
	if err := db.Raw(finalQuery, args...).Scan(&hotels).Error; err != nil {
		logger.Error(ctx, "Error fetching hotels (raw)", err.Error())
		return nil, 0, err
	}

	return hotels, total, nil
}

// buildBaseHotelQuery builds the core query structure reused across all filter functions
func (hr *HotelRepository) buildBaseHotelQuery(
	selectClause string,
	roomConditions []string,
	priceHaving string,
	hotelConditions []string,
	additionalJoins string,
	groupByClause string,
	orderByClause string,
) string {
	var queryBuilder strings.Builder

	// SELECT clause (custom per function)
	queryBuilder.WriteString(selectClause)
	queryBuilder.WriteString("\n\t\tFROM hotels h")

	// Subquery untuk minimum price (ALWAYS the same)
	queryBuilder.WriteString(`
		JOIN ( 
			SELECT rt.hotel_id, MIN(rp.price) AS min_price
			FROM room_types rt
			JOIN room_prices rp ON rt.id = rp.room_type_id
			JOIN bed_type_rooms btr ON btr.room_type_id = rt.id
			JOIN bed_types bt ON bt.id = btr.bed_type_id
			WHERE rp.is_show = true
	`)

	// Room conditions
	if len(roomConditions) > 0 {
		queryBuilder.WriteString("\n\t\t\tAND ")
		queryBuilder.WriteString(strings.Join(roomConditions, " AND "))
	}

	// GROUP BY untuk subquery
	queryBuilder.WriteString("\n\t\t\tGROUP BY rt.hotel_id")

	// HAVING untuk price filter
	if priceHaving != "" {
		queryBuilder.WriteString("\n\t\t\t")
		queryBuilder.WriteString(priceHaving)
	}

	// Close subquery
	queryBuilder.WriteString(`
		) mp ON mp.hotel_id = h.id
	`)

	// Additional JOINs (untuk GetFilterBedTypes, GetFilterTotalBedrooms)
	if additionalJoins != "" {
		queryBuilder.WriteString(additionalJoins)
	}

	// WHERE clause untuk hotel
	if len(hotelConditions) > 0 {
		queryBuilder.WriteString("\n\t\tWHERE h.deleted_at IS NULL AND ")
		queryBuilder.WriteString(strings.Join(hotelConditions, " AND "))
	}

	// GROUP BY (untuk aggregate queries)
	if groupByClause != "" {
		queryBuilder.WriteString("\n\t\t")
		queryBuilder.WriteString(groupByClause)
	}

	// ORDER BY
	if orderByClause != "" {
		queryBuilder.WriteString("\n\t\t")
		queryBuilder.WriteString(orderByClause)
	}

	return queryBuilder.String()
}
