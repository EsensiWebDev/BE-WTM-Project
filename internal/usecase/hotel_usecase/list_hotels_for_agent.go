package hotel_usecase

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) ListHotelsForAgent(ctx context.Context, req *hoteldto.ListHotelForAgentRequest) (*hoteldto.ListHotelForAgentResponse, error) {

	var rangeDateFrom, rangeDateTo time.Time
	var err error
	if req.RangeDateFrom != "" {
		rangeDateFrom, err = time.Parse(time.DateOnly, req.RangeDateFrom)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
			return nil, err
		}
	}
	if req.RangeDateTo != "" {
		rangeDateTo, err = time.Parse(time.DateOnly, req.RangeDateTo)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
			return nil, err
		}
	}

	filterHotel := filter.HotelFilterForAgent{
		Ratings:           req.Rating,
		BedTypeIDs:        req.BedTypeID,
		PaginationRequest: req.PaginationRequest,
		PriceMin:          req.RangePriceMin,
		PriceMax:          req.RangePriceMax,
		Cities:            req.District,
		TotalBedrooms:     req.TotalBedrooms,
		Province:          req.Province,
		DateFrom:          &rangeDateFrom,
		DateTo:            &rangeDateTo,
		PromoID:           uint(req.PromoID),
	}

	if req.TotalGuests > 0 && req.TotalRooms > 0 {
		filterHotel.MinGuest = req.TotalRooms / req.TotalGuests
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
					bucketName := fmt.Sprintf("%s-%s", constant.ConstHotel, constant.ConstPublic)
					photoUrl, err := hu.fileStorage.GetFile(ctx, bucketName, photo)
					if err != nil {
						logger.Error(ctx, "ListHotelsForAgent", err.Error())
					}
					respPhoto = photoUrl
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

		ratings = make([]entity.FilterRatingHotel, 0, 6)
		for i := 0; i <= 5; i++ {
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
