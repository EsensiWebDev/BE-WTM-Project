package hotel_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"

	"golang.org/x/sync/errgroup"
)

func (hu *HotelUsecase) ListHotelsForAgent(ctx context.Context, req *hoteldto.ListHotelForAgentRequest) (*hoteldto.ListHotelForAgentResponse, error) {

	var rangeDateFrom, rangeDateTo time.Time
	var err error
	resp := &hoteldto.ListHotelForAgentResponse{}
	if req.RangeDateFrom != "" {
		rangeDateFrom, err = time.Parse(time.DateOnly, req.RangeDateFrom)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
			return resp, nil
		}

		// Validasi: RangeDateFrom minimal hari ini
		today := time.Now().Truncate(24 * time.Hour)
		rangeDateFrom = rangeDateFrom.Truncate(24 * time.Hour)

		if rangeDateFrom.Before(today) {
			errMsg := "RangeDateFrom must be today or in the future"
			logger.Error(ctx, "ListHotelsForAgent", errMsg)
			return resp, nil
		}
	}

	if req.RangeDateTo != "" {
		rangeDateTo, err = time.Parse(time.DateOnly, req.RangeDateTo)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
			return resp, nil
		}

		// Validasi: RangeDateTo minimal hari ini
		today := time.Now().Truncate(24 * time.Hour)
		rangeDateTo = rangeDateTo.Truncate(24 * time.Hour)

		if rangeDateTo.Before(today) {
			errMsg := "RangeDateTo must be today or in the future"
			logger.Error(ctx, "ListHotelsForAgent", errMsg)
			return resp, nil
		}
	}

	// Validasi relasi antar tanggal
	if req.RangeDateFrom != "" && req.RangeDateTo != "" {
		// Validasi 1: RangeDateFrom harus sebelum RangeDateTo
		if rangeDateFrom.After(rangeDateTo) {
			errMsg := "RangeDateFrom must be before RangeDateTo"
			logger.Error(ctx, "ListHotelsForAgent", errMsg)
			return resp, nil
		}

		// Validasi 2: Tidak boleh sama (minimal 1 hari selisih)
		if rangeDateFrom.Equal(rangeDateTo) {
			errMsg := "RangeDateFrom and RangeDateTo cannot be the same date"
			logger.Error(ctx, "ListHotelsForAgent", errMsg)
			return resp, nil
		}
	}

	logger.Info(ctx, "Request ListHotelsForAgent", req)

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
		filterHotel.MinGuest = req.TotalGuests / req.TotalRooms
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
		logger.Info(ctx, "Request GetHotelsForAgent", filterHotel)
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
				Prices:   hotel.Prices,
				Currency: hotel.Currency,
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

	resp.Hotels = respHotels
	resp.FilterTotalRooms = totalRooms
	resp.FilterBedTypes = bedTypes
	resp.FilterRatings = ratings
	resp.FilterPricing = pricing
	resp.FilterDistricts = districts
	resp.Total = total

	return resp, nil

}
