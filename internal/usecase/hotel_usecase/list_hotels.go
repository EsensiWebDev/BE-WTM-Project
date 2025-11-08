package hotel_usecase

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListHotels(ctx context.Context, req *hoteldto.ListHotelRequest) (*hoteldto.ListHotelResponse, error) {
	filterHotel := filter.HotelFilter{
		IsAPI:             req.IsAPI,
		Region:            req.Region,
		PaginationRequest: req.PaginationRequest,
		StatusID:          req.StatusID,
	}

	hotels, total, err := hu.hotelRepo.GetHotels(ctx, filterHotel)
	if err != nil {
		logger.Error(ctx, "Error getting hotels", err.Error())
		return nil, err
	}

	response := &hoteldto.ListHotelResponse{
		Hotels: make([]hoteldto.ListHotel, 0, len(hotels)),
		Total:  total,
	}

	for _, hotelData := range hotels {
		data := hoteldto.ListHotel{
			ID:    hotelData.ID,
			Name:  hotelData.Name,
			IsAPI: hotelData.IsAPI,
		}

		if hotelData.AddrProvince != "" {
			data.Region = hotelData.AddrProvince
		}

		if hotelData.Email != "" {
			data.Email = hotelData.Email
		}

		if hotelData.StatusHotel != "" {
			data.Status = hotelData.StatusHotel
		}

		for _, roomType := range hotelData.RoomTypes {
			room := hoteldto.RoomTypeItem{
				Name:                  roomType.Name,
				Price:                 roomType.WithBreakfast.Price,
				PriceWithoutBreakfast: roomType.WithoutBreakfast.Price,
			}
			data.Rooms = append(data.Rooms, room)
		}

		response.Hotels = append(response.Hotels, data)
	}

	return response, nil
}
