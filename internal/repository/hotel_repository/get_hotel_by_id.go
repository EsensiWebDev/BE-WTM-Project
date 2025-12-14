package hotel_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"gorm.io/gorm"
)

func (hr *HotelRepository) GetHotelByID(ctx context.Context, hotelID uint, agentID uint) (*entity.Hotel, error) {
	db := hr.db.GetTx(ctx)

	var hotelModel model.Hotel

	query := db.WithContext(ctx).
		Preload("Status").
		Preload("HotelNearbyPlaces").
		Preload("HotelNearbyPlaces.NearbyPlace").
		Preload("Facilities").
		Preload("RoomTypes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("RoomTypes.BedTypes").
		Preload("RoomTypes.RoomTypeAdditionals").
		Preload("RoomTypes.RoomTypeAdditionals.RoomAdditional").
		Preload("RoomTypes.RoomPrices")

	if err := query.First(&hotelModel, hotelID).Error; err != nil {
		logger.Error(ctx, "Error fetching hotel by Id", err.Error())
		return nil, err
	}
	var promoAgent []model.Promo
	if agentID != 0 {
		queryPromo := db.WithContext(ctx).
			Model(&model.Promo{}).
			Preload("PromoRoomTypes").
			Where("is_active = ?", true).
			Where("id IN (?)",
				db.Table("detail_promo_groups").
					Select("promo_id").
					Where("promo_group_id IN (?)",
						db.Table("users").
							Select("promo_group_id").
							Where("id = ?", agentID),
					),
			)

		if err := queryPromo.Find(&promoAgent).Error; err != nil {
			logger.Error(ctx, "Error fetching promo by Id", err.Error())
			return nil, err
		}

	}

	var hotelEntity entity.Hotel
	if err := utils.CopyStrict(&hotelEntity, &hotelModel); err != nil {
		logger.Error(ctx, "Failed to copy hotel model to entity", err.Error())
		return nil, err
	}

	if hotelModel.SocialMedia != nil {
		if err := json.Unmarshal(hotelModel.SocialMedia, &hotelEntity.SocialMedia); err != nil {
			logger.Error(ctx, "Failed to unmarshal social media JSON", err.Error())
			return nil, err
		}
	}

	hotelEntity.StatusHotel = hotelModel.Status.Status
	for _, facility := range hotelModel.Facilities {
		hotelEntity.FacilityNames = append(hotelEntity.FacilityNames, facility.Name)
	}
	for _, nearbyPlace := range hotelModel.HotelNearbyPlaces {
		hotelEntity.NearbyPlaces = append(hotelEntity.NearbyPlaces, entity.NearbyPlace{
			ID:     nearbyPlace.NearbyPlaceID,
			Name:   nearbyPlace.NearbyPlace.Name,
			Radius: nearbyPlace.Radius,
		})
	}
	for i, roomType := range hotelModel.RoomTypes {
		for _, bedType := range roomType.BedTypes {
			hotelEntity.RoomTypes[i].BedTypeNames = append(hotelEntity.RoomTypes[i].BedTypeNames, bedType.Name)
		}
		for _, typeAdditional := range roomType.RoomTypeAdditionals {
			hotelEntity.RoomTypes[i].RoomAdditions = append(hotelEntity.RoomTypes[i].RoomAdditions, entity.CustomRoomAdditionalWithID{
				ID:         typeAdditional.ID,
				Name:       typeAdditional.RoomAdditional.Name,
				Category:   typeAdditional.Category,
				Price:      typeAdditional.Price,
				Pax:        typeAdditional.Pax,
				IsRequired: typeAdditional.IsRequired,
			})
		}

		for _, promo := range promoAgent {
			for _, promoRoomType := range promo.PromoRoomTypes {
				if promoRoomType.RoomTypeID == roomType.ID {
					var detailPromo entity.PromoDetail
					if err := json.Unmarshal(promo.Detail, &detailPromo); err != nil {
						logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
					}
					hotelEntity.RoomTypes[i].PromoRoomTypes = append(hotelEntity.RoomTypes[i].PromoRoomTypes, entity.PromoRoomTypes{
						ID:          promoRoomType.ID,
						PromoID:     promo.ID,
						RoomTypeID:  promoRoomType.RoomTypeID,
						TotalNights: promoRoomType.TotalNights,
						Promo: entity.Promo{
							ID:            promo.ID,
							ExternalID:    promo.ExternalID.ExternalID,
							Name:          promo.Name,
							StartDate:     promo.StartDate,
							EndDate:       promo.EndDate,
							Code:          promo.Code,
							Description:   promo.Description,
							PromoTypeID:   promo.PromoTypeID,
							Detail:        detailPromo,
							IsActive:      promo.IsActive,
							PromoTypeName: promo.PromoType.Name,
						},
					})
					break
				}
			}
		}
		for _, price := range roomType.RoomPrices {
			customBreakfast := entity.CustomBreakfastWithID{
				ID:     price.ID,
				Pax:    price.Pax,
				Price:  price.Price,
				IsShow: price.IsShow,
			}
			if price.IsBreakfast {
				hotelEntity.RoomTypes[i].WithBreakfast = customBreakfast
			} else if !price.IsBreakfast {
				hotelEntity.RoomTypes[i].WithoutBreakfast = customBreakfast
			}
		}
	}

	return &hotelEntity, nil
}
