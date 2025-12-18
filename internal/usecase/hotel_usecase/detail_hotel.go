package hotel_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) DetailHotel(ctx context.Context, hotelID uint) (*hoteldto.DetailHotelResponse, error) {
	hotel, err := hu.hotelRepo.GetHotelByID(ctx, hotelID, 0)
	if err != nil {
		logger.Error(ctx, "Error getting hotel by Id", err.Error())
		return nil, err
	}

	if hotel == nil {
		return nil, nil
	}

	bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, constant.ConstPublic)

	respHotel := &hoteldto.DetailHotelResponse{
		ID:                 hotel.ID,
		Name:               hotel.Name,
		Province:           hotel.AddrProvince,
		District:           hotel.AddrCity,
		SubDistrict:        hotel.AddrSubDistrict,
		Description:        hotel.Description,
		Rating:             hotel.Rating,
		Email:              hotel.Email,
		Facilities:         hotel.FacilityNames,
		NearbyPlace:        hotel.NearbyPlaces,
		CancellationPeriod: hotel.CancellationPeriod,
	}

	for _, photo := range hotel.Photos {
		photoUrl, err := hu.fileStorage.GetFile(ctx, bucketName, photo)
		if err != nil {
			logger.Error(ctx, "Error getting hotel photo", err.Error())
			return nil, fmt.Errorf("failed to get hotel photo: %s", err.Error())
		}
		respHotel.Photos = append(respHotel.Photos, photoUrl)
	}

	if hotel.CheckInHour != nil {
		respHotel.CheckInHour = hotel.CheckInHour.In(constant.AsiaJakarta).Format("15:04-07:00")
	}
	if hotel.CheckOutHour != nil {
		respHotel.CheckOutHour = hotel.CheckOutHour.In(constant.AsiaJakarta).Format("15:04-07:00")
	}

	var roomTypeList []hoteldto.DetailRoomType
	for _, rt := range hotel.RoomTypes {
		roomType := hoteldto.DetailRoomType{
			ID:                     rt.ID,
			Name:                   rt.Name,
			RoomSize:               rt.RoomSize,
			MaxOccupancy:           rt.MaxOccupancy,
			BedTypes:               rt.BedTypeNames,
			IsSmokingRoom:          rt.IsSmokingAllowed != nil && *rt.IsSmokingAllowed,
			Description:            rt.Description,
			BookingLimitPerBooking: rt.BookingLimitPerBooking,
		}

		for i, photo := range rt.Photos {
			photoUrl, err := hu.fileStorage.GetFile(ctx, bucketName, photo)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("Error getting room type photo index %d for room type ID %d", i, rt.ID), err.Error())
				return nil, fmt.Errorf("failed to get room type photo: %s", err.Error())
			}
			roomType.Photos = append(roomType.Photos, photoUrl)
		}

		roomType.WithoutBreakfast = entity.CustomBreakfast{
			Pax:    rt.WithoutBreakfast.Pax,
			Price:  rt.WithoutBreakfast.Price,  // DEPRECATED: Keep for backward compatibility
			Prices: rt.WithoutBreakfast.Prices, // NEW: Multi-currency prices
			IsShow: rt.WithoutBreakfast.IsShow,
		}

		roomType.WithBreakfast = entity.CustomBreakfast{
			Pax:    rt.WithBreakfast.Pax,
			Price:  rt.WithBreakfast.Price,  // DEPRECATED: Keep for backward compatibility
			Prices: rt.WithBreakfast.Prices, // NEW: Multi-currency prices
			IsShow: rt.WithBreakfast.IsShow,
		}

		for _, addition := range rt.RoomAdditions {
			roomType.Additional = append(roomType.Additional, entity.CustomRoomAdditionalWithID{
				ID:         addition.ID,
				Name:       addition.Name,
				Category:   addition.Category,
				Price:      addition.Price,  // DEPRECATED: Keep for backward compatibility
				Prices:     addition.Prices, // NEW: Multi-currency prices
				Pax:        addition.Pax,
				IsRequired: addition.IsRequired,
			})
		}

		for _, pref := range rt.OtherPreferences {
			roomType.OtherPreferences = append(roomType.OtherPreferences, entity.CustomOtherPreferenceWithID{
				ID:   pref.ID,
				Name: pref.Name,
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
