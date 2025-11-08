package hotel_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) DetailHotel(ctx context.Context, hotelID uint) (*hoteldto.DetailHotelResponse, error) {
	hotel, err := hu.hotelRepo.GetHotelByID(ctx, hotelID, constant.RoleAdmin)
	if err != nil {
		logger.Error(ctx, "Error getting hotel by Id", err.Error())
		return nil, err
	}

	if hotel == nil {
		return nil, nil
	}

	respHotel := &hoteldto.DetailHotelResponse{
		ID:                 hotel.ID,
		Name:               hotel.Name,
		Province:           hotel.AddrProvince,
		District:           hotel.AddrCity,
		SubDistrict:        hotel.AddrSubDistrict,
		Description:        hotel.Description,
		Photos:             hotel.Photos,
		Rating:             hotel.Rating,
		Email:              hotel.Email,
		Facilities:         hotel.FacilityNames,
		NearbyPlace:        hotel.NearbyPlaces,
		CancellationPeriod: hotel.CancellationPeriod,
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
			ID:            rt.ID,
			Name:          rt.Name,
			RoomSize:      rt.RoomSize,
			MaxOccupancy:  rt.MaxOccupancy,
			BedTypes:      rt.BedTypeNames,
			IsSmokingRoom: rt.IsSmokingAllowed != nil && *rt.IsSmokingAllowed,
			Description:   rt.Description,
			Photos:        rt.Photos,
		}

		roomType.WithoutBreakfast = entity.CustomBreakfast{
			Pax:    rt.WithoutBreakfast.Pax,
			Price:  rt.WithoutBreakfast.Price,
			IsShow: rt.WithoutBreakfast.IsShow,
		}

		roomType.WithBreakfast = entity.CustomBreakfast{
			Pax:    rt.WithBreakfast.Pax,
			Price:  rt.WithBreakfast.Price,
			IsShow: rt.WithBreakfast.IsShow,
		}

		for _, addition := range rt.RoomAdditions {
			roomType.Additional = append(roomType.Additional, entity.CustomRoomAdditionalWithID{
				ID:    addition.ID,
				Name:  addition.Name,
				Price: addition.Price,
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
