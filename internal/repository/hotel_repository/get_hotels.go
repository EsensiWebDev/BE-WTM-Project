package hotel_repository

import (
	"context"
	"gorm.io/gorm/clause"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

var validateColumnSort = map[string]bool{}

func (hr *HotelRepository) GetHotels(ctx context.Context, filter filter.HotelFilter) ([]entity.Hotel, int64, error) {
	db := hr.db.GetTx(ctx)
	// Select default fields
	selectFields := []string{"id", "name", "addr_province", "email", "status_id", "is_api"}

	// Initialize query with default fields
	query := db.WithContext(ctx).
		Model(&model.Hotel{}).
		Preload("RoomTypes").
		Preload("RoomTypes.RoomPrices").
		Select(selectFields)

	// Apply filters
	if filter.IsAPI != nil {
		isAPI := *filter.IsAPI
		query = query.Where("is_api = ?", isAPI)
	}

	if len(filter.Region) > 0 {
		query = query.Where("addr_province IN ?", filter.Region)
	}

	if filter.StatusID > 0 {
		query = query.Where("status_id = ?", filter.StatusID)
	}

	// Search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(name) ILIKE ? ", "%"+safeSearch+"%")
	}

	// Count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting hotels", err.Error())
		return nil, total, err
	}

	// Pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// Order
	if filter.Sort != "" {
		if validateColumnSort[filter.Sort] {
			var desc bool
			if strings.TrimSpace(strings.ToLower(filter.Dir)) == "asc" {
				desc = false
			} else {
				desc = true
			}
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: filter.Sort}, Desc: desc})
		}
	} else {
		query = query.Order("id desc")
	}

	// Execute
	var hotels []model.Hotel
	if err := query.Find(&hotels).Error; err != nil {
		logger.Error(ctx, "Error fetching hotels", err.Error())
		return nil, total, err
	}

	// Mapping
	var hotelEntities []entity.Hotel
	for _, hotel := range hotels {
		var hotelEntity entity.Hotel
		if err := utils.CopyPatch(&hotelEntity, &hotel); err != nil {
			logger.Error(ctx, "Failed to copy hotel model to entity", err.Error())
			return nil, total, err
		}
		for i, roomType := range hotel.RoomTypes {
			for _, price := range roomType.RoomPrices {
				if price.IsBreakfast {
					hotelEntity.RoomTypes[i].WithBreakfast = entity.CustomBreakfastWithID{
						ID:     price.ID,
						Price:  price.Price,
						Pax:    price.Pax,
						IsShow: price.IsShow,
					}
				} else {
					hotelEntity.RoomTypes[i].WithoutBreakfast = entity.CustomBreakfastWithID{
						ID:     price.ID,
						Price:  price.Price,
						Pax:    price.Pax,
						IsShow: price.IsShow,
					}
				}
			}

		}

		hotelEntity.StatusHotel = constant.MapStatusHotel[int(hotel.StatusID)]
		hotelEntities = append(hotelEntities, hotelEntity)
	}

	return hotelEntities, total, nil

}
