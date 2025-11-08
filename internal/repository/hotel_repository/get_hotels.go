package hotel_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetHotels(ctx context.Context, filter filter.HotelFilter) ([]entity.Hotel, int64, error) {
	db := hr.db.GetTx(ctx)
	// Select default fields
	selectFields := []string{"id", "name", "addr_province", "email", "status_id", "is_api"}

	// Initialize query with default fields
	query := db.WithContext(ctx).
		Model(&model.Hotel{}).
		Preload("RoomTypes").
		Preload("RoomTypes.RoomPrices").
		Select(selectFields).Debug()

	// Apply filters
	if filter.IsAPI != nil {
		isAPI := *filter.IsAPI
		query = query.Where("is_api = ?", isAPI)
	}

	if filter.Region != "" {
		query = query.Where("addr_province = ?", filter.Region)
	}

	if filter.StatusID > 0 {
		query = query.Where("status_id = ?", filter.StatusID)
	}

	// Search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(name) ILIKE ? ESCAPE '\\'", "%"+safeSearch+"%")
	}

	// Count
	var total int64
	if err := query.Debug().Count(&total).Error; err != nil {
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
