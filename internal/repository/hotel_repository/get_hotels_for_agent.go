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

func (hr *HotelRepository) GetHotelsForAgent(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.CustomHotel, int64, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var hotelConditions []string
	var roomConditions []string
	var priceHaving string

	// 🔍 Filter kota
	if len(filter.Cities) > 0 {
		hotelConditions = append(hotelConditions, "h.addr_city IN (?)")
		args = append(args, filter.Cities)
	}

	// 🔍 Filter rating
	if len(filter.Ratings) > 0 {
		hotelConditions = append(hotelConditions, "h.rating IN (?)")
		args = append(args, filter.Ratings)
	}

	// 🔍 Filter nama hotel
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		hotelConditions = append(hotelConditions, "LOWER(h.name) ILIKE ? ESCAPE '\\'")
		args = append(args, "%"+safeSearch+"%")
	}

	// 🔍 Filter bed type dan total unit
	if len(filter.BedTypeIDs) > 0 {
		roomConditions = append(roomConditions, "bt.id IN (?)")
		args = append(args, filter.BedTypeIDs)
	}
	if len(filter.TotalBedrooms) > 0 {
		roomConditions = append(roomConditions, "rt.total_unit IN (?)")
		args = append(args, filter.TotalBedrooms)
	}

	// 🔍 Filter harga
	if filter.PriceMin != nil && filter.PriceMax != nil {
		priceHaving = "HAVING MIN(rp.price) BETWEEN ? AND ?"
		args = append(args, *filter.PriceMin, *filter.PriceMax)
	} else if filter.PriceMin != nil {
		priceHaving = "HAVING MIN(rp.price) > ?"
		args = append(args, *filter.PriceMin)
	} else if filter.PriceMax != nil {
		priceHaving = "HAVING MIN(rp.price) < ?"
		args = append(args, *filter.PriceMax)
	}

	// 🧩 WHERE clause hotel
	hotelWhere := ""
	if len(hotelConditions) > 0 {
		hotelWhere = "WHERE " + strings.Join(hotelConditions, " AND ")
	}

	// 🧩 WHERE clause room
	roomWhere := ""
	if len(roomConditions) > 0 {
		roomWhere = "AND " + strings.Join(roomConditions, " AND ")
	}

	// 🧠 Final raw SQL
	rawQuery := fmt.Sprintf(`
        SELECT h.id, h.name, h.addr_province, h.addr_city, h.addr_sub_district, h.photos, h.rating, mp.min_price
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
            %s
        ) mp ON mp.hotel_id = h.id
        %s
    `, roomWhere, priceHaving, hotelWhere)

	// 📦 Pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		rawQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)
	}

	// 🔍 Execute query
	var hotels []entity.CustomHotel
	if err := db.Raw(rawQuery, args...).Scan(&hotels).Error; err != nil {
		logger.Error(ctx, "Error fetching hotels (raw)", err.Error())
		return nil, 0, err
	}

	// 🔢 Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS sub", rawQuery)
	var total int64
	if err := db.Raw(countQuery, args...).Debug().Scan(&total).Error; err != nil {
		logger.Error(ctx, "Error counting hotels (raw)", err.Error())
		return nil, 0, err
	}

	return hotels, total, nil
}
