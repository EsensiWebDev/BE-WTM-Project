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

func (hr *HotelRepository) GetFilterRatings(ctx context.Context, filter filter.HotelFilterForAgent) ([]entity.FilterRatingHotel, error) {
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

	// ⚠️ TIDAK include Ratings (karena ini fungsi untuk get available ratings)

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
		`SELECT h.rating, COUNT(DISTINCT h.id) AS count`,
		roomConditions,
		priceHaving,
		hotelConditions,
		"", // no additional joins
		"GROUP BY h.rating",
		"ORDER BY h.rating ASC",
	)

	// 🔍 Execute query
	var ratings []entity.FilterRatingHotel
	if err := db.Raw(query, args...).Scan(&ratings).Error; err != nil {
		logger.Error(ctx, "Error fetching filter ratings (raw)", err.Error())
		return nil, err
	}

	if len(ratings) == 0 {
		logger.Info(ctx, "No ratings found for the given filter criteria")
		return []entity.FilterRatingHotel{}, nil
	}

	return ratings, nil
}
