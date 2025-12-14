package hotel_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) DetailHotelForAgent(ctx context.Context, hotelID uint) (*hoteldto.DetailHotelForAgentResponse, error) {
	userCtx, err := hu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return nil, err
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return nil, err
	}

	agentID := userCtx.ID

	hotel, err := hu.hotelRepo.GetHotelByID(ctx, hotelID, agentID)
	if err != nil {
		logger.Error(ctx, "Error getting hotel by Id", err.Error())
		return nil, err
	}

	if hotel == nil {
		return nil, nil
	}

	respHotel := &hoteldto.DetailHotelForAgentResponse{
		ID:                 hotel.ID,
		Name:               hotel.Name,
		Province:           hotel.AddrProvince,
		District:           hotel.AddrCity,
		SubDistrict:        hotel.AddrSubDistrict,
		Description:        hotel.Description,
		Rating:             hotel.Rating,
		Email:              hotel.Email,
		Facilities:         hotel.FacilityNames,
		CancellationPeriod: hotel.CancellationPeriod,
	}

	bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, constant.ConstPublic)
	for _, photo := range hotel.Photos {
		photoUrl, err := hu.fileStorage.GetFile(ctx, bucketName, photo)
		if err != nil {
			logger.Error(ctx, "Error getting user profile photo", err.Error())
			return nil, fmt.Errorf("failed to get user profile photo: %s", err.Error())
		}
		respHotel.Photos = append(respHotel.Photos, photoUrl)
	}

	var nearbyPlaces []hoteldto.NearbyPlaceForAgent
	for _, nearbyPlace := range hotel.NearbyPlaces {
		nearbyPlaces = append(nearbyPlaces, hoteldto.NearbyPlaceForAgent{
			Name:   nearbyPlace.Name,
			Radius: nearbyPlace.Radius,
		})
	}
	respHotel.NearbyPlace = nearbyPlaces

	if hotel.CheckInHour != nil {
		respHotel.CheckInHour = hotel.CheckInHour.In(constant.AsiaJakarta).Format("15:04-07:00")
	}
	if hotel.CheckOutHour != nil {
		respHotel.CheckOutHour = hotel.CheckOutHour.In(constant.AsiaJakarta).Format("15:04-07:00")
	}

	var roomTypeList []hoteldto.DetailRoomTypeForAgent
	for _, rt := range hotel.RoomTypes {
		roomType := hoteldto.DetailRoomTypeForAgent{
			Name:          rt.Name,
			RoomSize:      rt.RoomSize,
			MaxOccupancy:  rt.MaxOccupancy,
			BedTypes:      rt.BedTypeNames,
			IsSmokingRoom: rt.IsSmokingAllowed != nil && *rt.IsSmokingAllowed,
			Description:   rt.Description,
		}
		for _, photo := range rt.Photos {
			photoUrl, err := hu.fileStorage.GetFile(ctx, bucketName, photo)
			if err != nil {
				logger.Error(ctx, "Error getting user profile photo", err.Error())
				return nil, fmt.Errorf("failed to get user profile photo: %s", err.Error())
			}
			roomType.Photos = append(roomType.Photos, photoUrl)
		}

		roomType.WithoutBreakfast = entity.CustomBreakfastWithID{
			ID:     rt.WithoutBreakfast.ID,
			Pax:    rt.WithoutBreakfast.Pax,
			Price:  rt.WithoutBreakfast.Price,
			IsShow: rt.WithoutBreakfast.IsShow,
		}

		roomType.WithBreakfast = entity.CustomBreakfastWithID{
			ID:     rt.WithBreakfast.ID,
			Pax:    rt.WithBreakfast.Pax,
			Price:  rt.WithBreakfast.Price,
			IsShow: rt.WithBreakfast.IsShow,
		}

		var promos []hoteldto.PromoDetailRoom
		for _, prt := range rt.PromoRoomTypes {
			if prt.Promo.IsActive {
				var priceWithBreakfast, priceWithoutBreakfast float64
				var notes string
				if prt.Promo.Detail.DiscountPercentage > 0 {
					priceWithBreakfast = (100 - prt.Promo.Detail.DiscountPercentage) / 100 * roomType.WithBreakfast.Price
					priceWithoutBreakfast = (100 - prt.Promo.Detail.DiscountPercentage) / 100 * roomType.WithoutBreakfast.Price
				} else if prt.Promo.Detail.FixedPrice > 0 {
					priceWithBreakfast = prt.Promo.Detail.FixedPrice
					priceWithoutBreakfast = prt.Promo.Detail.FixedPrice
				} else if prt.Promo.Detail.BenefitNote != "" {
					notes = prt.Promo.Detail.BenefitNote
				}
				promos = append(promos, hoteldto.PromoDetailRoom{
					PromoID:               prt.Promo.ID,
					TotalNights:           prt.TotalNights,
					Description:           prt.Promo.Description,
					CodePromo:             prt.Promo.Code,
					PriceWithBreakfast:    priceWithBreakfast,
					PriceWithoutBreakfast: priceWithoutBreakfast,
					OtherNotes:            notes,
				})
			}
		}
		roomType.Promos = promos

		for _, addition := range rt.RoomAdditions {
			roomType.Additional = append(roomType.Additional, entity.CustomRoomAdditionalWithID{
				ID:         addition.ID,
				Name:       addition.Name,
				Category:   addition.Category,
				Price:      addition.Price,
				Pax:        addition.Pax,
				IsRequired: addition.IsRequired,
			})
		}

		roomTypeList = append(roomTypeList, roomType)
	}

	respHotel.RoomType = roomTypeList

	var sosialMediaList []hoteldto.SocialMedia
	for typeSM, url := range hotel.SocialMedia {
		sosialMediaList = append(sosialMediaList, hoteldto.SocialMedia{
			Platform: typeSM,
			Link:     url,
		})
	}

	respHotel.SocialMedia = sosialMediaList

	return respHotel, nil
}
