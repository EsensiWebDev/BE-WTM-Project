package hotel_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) GetHotelByID(ctx context.Context, hotelID uint, scope string) (*entity.Hotel, error) {
	db := hr.db.GetTx(ctx)

	var hotelModel model.Hotel

	query := db.WithContext(ctx).
		Preload("Status").
		Preload("HotelNearbyPlaces").
		Preload("HotelNearbyPlaces.NearbyPlace").
		Preload("Facilities").
		Preload("RoomTypes").
		Preload("RoomTypes.BedTypes").
		Preload("RoomTypes.RoomTypeAdditionals").
		Preload("RoomTypes.RoomTypeAdditionals.RoomAdditional").
		Preload("RoomTypes.RoomPrices")

	if scope == constant.RoleAgent {
		query.Preload("RoomTypes.PromoRoomTypes").Preload("RoomTypes.PromoRoomTypes.Promo")
	}

	if err := query.First(&hotelModel, hotelID).Error; err != nil {
		logger.Error(ctx, "Error fetching hotel by Id", err.Error())
		return nil, err
	}

	var hotelEntity entity.Hotel
	if err := utils.CopyStrict(&hotelEntity, &hotelModel); err != nil {
		logger.Error(ctx, "Failed to copy hotel model to entity", err.Error())
		return nil, err
	}

	if err := json.Unmarshal(hotelModel.SocialMedia, &hotelEntity.SocialMedia); err != nil {
		logger.Error(ctx, "Failed to unmarshal social media JSON", err.Error())
		return nil, err
	}

	hotelEntity.StatusHotel = hotelModel.Status.Status
	for _, facility := range hotelModel.Facilities {
		hotelEntity.FacilityNames = append(hotelEntity.FacilityNames, facility.Name)
	}
	for i, roomType := range hotelModel.RoomTypes {
		for _, bedType := range roomType.BedTypes {
			hotelEntity.RoomTypes[i].BedTypeNames = append(hotelEntity.RoomTypes[i].BedTypeNames, bedType.Name)
		}
		for _, typeAdditional := range roomType.RoomTypeAdditionals {
			hotelEntity.RoomTypes[i].RoomAdditions = append(hotelEntity.RoomTypes[i].RoomAdditions, entity.CustomRoomAdditionalWithID{
				ID:    typeAdditional.ID,
				Name:  typeAdditional.RoomAdditional.Name,
				Price: typeAdditional.Price,
			})
		}
		for i2, promoRoomType := range roomType.PromoRoomTypes {
			if len(promoRoomType.Promo.Detail) > 0 && promoRoomType.Promo.IsActive {
				var detailPromo entity.PromoDetail
				if err := json.Unmarshal(promoRoomType.Promo.Detail, &detailPromo); err != nil {
					logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
				}
				hotelEntity.RoomTypes[i].PromoRoomTypes[i2].Promo.Detail = detailPromo
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
		for _, nearbyPlace := range hotelModel.HotelNearbyPlaces {
			hotelEntity.NearbyPlaces = append(hotelEntity.NearbyPlaces, entity.NearbyPlace{
				ID:     nearbyPlace.NearbyPlaceID,
				Name:   nearbyPlace.NearbyPlace.Name,
				Radius: nearbyPlace.Radius,
			})
		}
	}

	return &hotelEntity, nil
}
