package hotel_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) GetFilterTotalBedrooms(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterTotalBedroom, error) {
	db := hr.db.GetTx(ctx)

	var args []interface{}
	var hotelConditions []string
	var roomConditions []string
	var priceHaving string

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

	// 🔍 Filter bed type
	if len(filter.BedTypeIDs) > 0 {
		roomConditions = append(roomConditions, "bt.id IN ?")
		args = append(args, filter.BedTypeIDs)
	}

	// 🔍 Filter availability
	if filter.DateFrom != nil && filter.DateTo != nil {
		roomConditions = append(roomConditions,
			"NOT EXISTS (SELECT 1 FROM room_unavailables ru WHERE ru.room_type_id = rt.id AND ru.date BETWEEN ? AND ?)")
		args = append(args, *filter.DateFrom, *filter.DateTo)
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

	// 🔍 Filter status hotel
	hotelConditions = append(hotelConditions, "h.status_id = ?")
	args = append(args, constant.StatusHotelApprovedID)

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
        SELECT rt.total_unit AS total_bed_rooms, COUNT(DISTINCT h.id) AS count
        FROM hotels h
        JOIN (
            SELECT rt.hotel_id, MIN(rp.price) AS min_price
            FROM room_types rt
            JOIN room_prices rp ON rt.id = rp.room_type_id
            WHERE rp.is_show = true
            GROUP BY rt.hotel_id
            %s
        ) mp ON mp.hotel_id = h.id
        JOIN (
            SELECT rt.hotel_id, rt.total_unit
            FROM room_types rt
            JOIN bed_type_rooms btr ON btr.room_type_id = rt.id
            JOIN bed_types bt ON bt.id = btr.bed_type_id
            %s
        ) rt ON rt.hotel_id = h.id
        %s
        GROUP BY rt.total_unit
    `, priceHaving, roomWhere, hotelWhere)

	var totalRooms []entity.FilterTotalBedroom
	if err := db.Raw(rawQuery, args...).Debug().Scan(&totalRooms).Error; err != nil {
		logger.Error(ctx, "Error fetching total bedrooms (raw)", err.Error())
		return nil, err
	}

	if len(totalRooms) == 0 {
		logger.Info(ctx, "No rooms found for the given filter criteria")
		return nil, nil
	}

	return totalRooms, nil
}
