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

func (hr *HotelRepository) GetFilterPricing(ctx context.Context, filter filter.HotelFilterForAgent) (*entity.FilterRangePrice, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var conditions []string

	// 🔍 Filter kota
	if len(filter.Cities) > 0 {
		conditions = append(conditions, "h.addr_city IN (?)")
		args = append(args, filter.Cities)
	}

	// 🔍 Filter rating
	if len(filter.Ratings) > 0 {
		conditions = append(conditions, "h.rating IN (?)")
		args = append(args, filter.Ratings)
	}

	// 🔍 Filter nama hotel
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		conditions = append(conditions, "LOWER(h.name) ILIKE ? ESCAPE '\\'")
		args = append(args, "%"+safeSearch+"%")
	}

	// 🔍 Filter bed type dan total unit
	var roomConditions []string
	if len(filter.BedTypeIDs) > 0 {
		roomConditions = append(roomConditions, "bt.id IN (?)")
		args = append(args, filter.BedTypeIDs)
	}
	if len(filter.TotalBedrooms) > 0 {
		roomConditions = append(roomConditions, "rt.total_unit IN (?)")
		args = append(args, filter.TotalBedrooms)
	}

	// 🧩 WHERE clause hotel
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 🧩 WHERE clause room type
	roomFilterClause := ""
	if len(roomConditions) > 0 {
		roomFilterClause = "AND " + strings.Join(roomConditions, " AND ")
	}

	// 🧠 Final raw SQL
	rawQuery := fmt.Sprintf(`
        SELECT MIN(mp.min_price) AS min_price, MAX(mp.min_price) AS max_price
        FROM hotels h
        JOIN (
            SELECT rt.hotel_id, MIN(rp.price) AS min_price
            FROM room_types rt
            JOIN room_prices rp ON rt.id = rp.room_type_id
            JOIN bed_type_rooms btr ON btr.room_type_id = rt.id
            JOIN bed_types bt ON bt.id = btr.bed_type_id
            WHERE rp.is_show = true
            %s
            GROUP BY rt.hotel_id
        ) mp ON mp.hotel_id = h.id
        %s
    `, roomFilterClause, whereClause)

	var result entity.FilterRangePrice
	if err := db.Raw(rawQuery, args...).Debug().Scan(&result).Error; err != nil {
		logger.Error(ctx, "Error fetching range price (raw)", err.Error())
		return nil, err
	}

	return &result, nil
}
