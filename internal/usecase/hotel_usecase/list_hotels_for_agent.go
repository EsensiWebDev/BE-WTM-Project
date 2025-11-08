package hotel_usecase

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListHotelsForAgent(ctx context.Context, req *hoteldto.ListHotelForAgentRequest) (*hoteldto.ListHotelForAgentResponse, error) {
	filterHotel := filter.HotelFilterForAgent{
		Ratings:           req.Rating,
		BedTypeIDs:        req.BedTypeID,
		PaginationRequest: req.PaginationRequest,
		PriceMin:          req.RangePriceMin,
		PriceMax:          req.RangePriceMax,
		Cities:            req.District,
		TotalBedrooms:     req.TotalRooms,
		Province:          req.Province,
		DateFrom:          req.RangeDateFrom,
		DateTo:            req.RangeDateTo,
		TotalGuest:        req.TotalGuests,
	}

	filter.Clean(&filterHotel)

	var (
		hotels     []entity.CustomHotel
		respHotels []hoteldto.ListHotelForAgent
		total      int64
		districts  []string
		pricing    *entity.FilterRangePrice
		ratings    []entity.FilterRatingHotel
		bedTypes   []entity.FilterBedTypeHotel
		totalRooms []entity.FilterTotalBedroom
	)

	// ⛳ Gunakan errgroup
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err error
		hotels, total, err = hu.hotelRepo.GetHotelsForAgent(egCtx, filterHotel)
		respHotels = make([]hoteldto.ListHotelForAgent, 0, len(hotels))
		for _, hotel := range hotels {
			var respPhoto string
			for _, photo := range hotel.Photos {
				if photo != "" {
					respPhoto = photo
					break
				}
			}
			respHotels = append(respHotels, hoteldto.ListHotelForAgent{
				ID:       hotel.ID,
				Name:     hotel.Name,
				Address:  fmt.Sprintf("%s, %s, %s", hotel.AddrSubDistrict, hotel.AddrCity, hotel.AddrProvince),
				MinPrice: hotel.MinPrice,
				Photo:    respPhoto,
				Rating:   hotel.Rating,
			})
		}
		return err
	})
	eg.Go(func() error {
		var err error
		districts, err = hu.hotelRepo.GetFilterDistricts(egCtx, filterHotel)
		return err
	})
	eg.Go(func() error {
		var err error
		pricing, err = hu.hotelRepo.GetFilterPricing(egCtx, filterHotel)
		return err
	})
	eg.Go(func() error {
		dataRatings, err := hu.hotelRepo.GetFilterRatings(egCtx, filterHotel)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
			return err
		}

		ratings = make([]entity.FilterRatingHotel, 0, 5)
		for i := 1; i <= 5; i++ {
			rate := entity.FilterRatingHotel{
				Rating: i,
			}

			for _, rating := range dataRatings {
				if rating.Rating == i {
					rate.Count = rating.Count
					break
				}
			}

			ratings = append(ratings, rate)
		}

		return nil
	})
	eg.Go(func() error {
		var err error
		bedTypes, err = hu.hotelRepo.GetFilterBedTypes(egCtx, filterHotel)
		return err
	})
	eg.Go(func() error {
		var err error
		totalRooms, err = hu.hotelRepo.GetFilterTotalBedrooms(egCtx, filterHotel)
		return err
	})

	if err := eg.Wait(); err != nil {
		logger.Error(ctx, "ListHotelsForAgent", err.Error())
		return nil, err
	}

	resp := &hoteldto.ListHotelForAgentResponse{
		Hotels:           respHotels,
		FilterTotalRooms: totalRooms,
		FilterBedTypes:   bedTypes,
		FilterRatings:    ratings,
		FilterPricing:    pricing,
		FilterDistricts:  districts,
		Total:            total,
	}

	return resp, nil

}
