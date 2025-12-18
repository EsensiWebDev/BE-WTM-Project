package hotel_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/currency"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/lib/pq"
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
	// Use a scan struct without Prices field to avoid GORM scanning errors
	// Photos is scanned as string (PostgreSQL array format) and parsed
	type HotelScan struct {
		ID              uint
		Name            string
		AddrSubDistrict string
		AddrCity        string
		AddrProvince    string
		Photos          string // Scan as string, will parse PostgreSQL array format
		Rating          int
		MinPrice        float64
	}
	var hotelScans []HotelScan
	if err := db.Raw(finalQuery, args...).Scan(&hotelScans).Error; err != nil {
		logger.Error(ctx, "Error fetching hotels (raw)", err.Error())
		return nil, 0, err
	}

	// Convert scan results to CustomHotel entities
	hotels := make([]entity.CustomHotel, len(hotelScans))
	for i, scan := range hotelScans {
		// Parse PostgreSQL array format string to []string
		// Format: {value1,value2,value3} or NULL
		var photos pq.StringArray
		if scan.Photos != "" {
			// Remove curly braces and split by comma
			photosStr := strings.Trim(scan.Photos, "{}")
			if photosStr != "" {
				photoList := strings.Split(photosStr, ",")
				photos = make(pq.StringArray, len(photoList))
				for j, photo := range photoList {
					photos[j] = strings.TrimSpace(photo)
				}
			}
		}
		if photos == nil {
			photos = pq.StringArray{} // Initialize empty array
		}

		hotels[i] = entity.CustomHotel{
			ID:              scan.ID,
			Name:            scan.Name,
			AddrSubDistrict: scan.AddrSubDistrict,
			AddrCity:        scan.AddrCity,
			AddrProvince:    scan.AddrProvince,
			Photos:          photos,
			Rating:          scan.Rating,
			MinPrice:        scan.MinPrice,
			Prices:          make(map[string]float64), // Initialize empty, will be populated from room_prices
		}
	}

	// Extract prices from room_prices for each hotel
	if len(hotels) > 0 {
		hotelIDs := make([]uint, len(hotels))
		for i, hotel := range hotels {
			hotelIDs[i] = hotel.ID
		}

		// Query all room_prices for these hotels
		type RoomPriceRow struct {
			HotelID uint
			Prices  []byte
		}
		var roomPrices []RoomPriceRow
		err := db.Table("room_prices rp").
			Select("rt.hotel_id, rp.prices").
			Joins("JOIN room_types rt ON rt.id = rp.room_type_id").
			Where("rt.hotel_id IN ? AND rp.is_show = true AND rp.prices IS NOT NULL AND rp.prices != '{}'::jsonb", hotelIDs).
			Scan(&roomPrices).Error

		if err == nil {
			// Group prices by hotel_id and find minimum for each currency
			hotelPricesMap := make(map[uint]map[string]float64)
			for _, rp := range roomPrices {
				prices, err := currency.JSONToPrices(rp.Prices)
				if err != nil {
					logger.Error(ctx, "Failed to parse prices JSONB", err.Error())
					continue
				}

				if hotelPricesMap[rp.HotelID] == nil {
					hotelPricesMap[rp.HotelID] = make(map[string]float64)
				}

				// Find minimum price for each currency
				for curr, price := range prices {
					if existingPrice, exists := hotelPricesMap[rp.HotelID][curr]; !exists || price < existingPrice {
						hotelPricesMap[rp.HotelID][curr] = price
					}
				}
			}

			// Assign prices to hotels
			for i := range hotels {
				if prices, exists := hotelPricesMap[hotels[i].ID]; exists && len(prices) > 0 {
					hotels[i].Prices = prices
				}
			}
		} else {
			logger.Error(ctx, "Failed to fetch prices for hotels", err.Error())
		}
	}

	// Filter hotels that have prices in the agent's currency
	// Only filter if currency is specified and not empty
	if filter.Currency != "" {
		normalizedCurrency := currency.NormalizeCurrencyCode(filter.Currency)
		originalCount := len(hotels)
		filteredHotels := make([]entity.CustomHotel, 0, len(hotels))

		for _, hotel := range hotels {
			// Only include hotels that have a Prices map AND contain the agent's currency with valid price > 0
			// Hotels without Prices map or without the agent's currency are excluded
			hasValidPrice := false

			if len(hotel.Prices) > 0 {
				if price, exists := hotel.Prices[normalizedCurrency]; exists && price > 0 {
					hasValidPrice = true
				} else {
					// Log why hotel is being excluded
					availableCurrencies := make([]string, 0, len(hotel.Prices))
					for curr := range hotel.Prices {
						availableCurrencies = append(availableCurrencies, curr)
					}
					logger.Info(ctx, "Excluding hotel - missing currency in prices", map[string]interface{}{
						"hotel_id":             hotel.ID,
						"hotel_name":           hotel.Name,
						"required_currency":    normalizedCurrency,
						"available_currencies": availableCurrencies,
						"has_currency":         exists,
						"price_value":          price,
					})
				}
			} else {
				// Hotel has no Prices map
				logger.Info(ctx, "Excluding hotel - no prices map", map[string]interface{}{
					"hotel_id":          hotel.ID,
					"hotel_name":        hotel.Name,
					"required_currency": normalizedCurrency,
				})
			}

			if hasValidPrice {
				filteredHotels = append(filteredHotels, hotel)
			}
		}

		hotels = filteredHotels
		// Update total count to reflect filtered results
		total = int64(len(hotels))
		logger.Info(ctx, "Filtered hotels by currency", map[string]interface{}{
			"currency":       normalizedCurrency,
			"original_count": originalCount,
			"filtered_count": len(hotels),
			"excluded_count": originalCount - len(hotels),
		})
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
