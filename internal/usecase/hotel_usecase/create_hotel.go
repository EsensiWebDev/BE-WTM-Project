package hotel_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hu *HotelUsecase) CreateHotel(ctx context.Context, req *hoteldto.CreateHotelRequest) (*hoteldto.CreateHotelResponse, error) {
	var hotelID uint
	err := hu.dbTransaction.WithTransaction(ctx, func(txCtx context.Context) error {
		var nearbyPlaces []hoteldto.NearbyPlace
		if req.NearbyPlaces != "" {
			if err := json.Unmarshal([]byte(req.NearbyPlaces), &nearbyPlaces); err != nil {
				logger.Error(ctx, "Failed to unmarshal CreateHotelRequest-NearbyPlaces", err.Error())
				return err
			}
		}

		var socialMedias []hoteldto.SocialMedia
		if req.SocialMedias != "" {
			if err := json.Unmarshal([]byte(req.SocialMedias), &socialMedias); err != nil {
				logger.Error(ctx, "Failed to unmarshal CreateHotelRequest-SocialMedias", err.Error())
				return err
			}
		}

		checkInHour, err := utils.ParseHourString(hu.config.DefaultCheckInHour)
		if err != nil {
			logger.Error(ctx, "Failed to parse checkInHour", err.Error())
		}

		checkOutHour, err := utils.ParseHourString(hu.config.DefaultCheckOutHour)
		if err != nil {
			logger.Error(ctx, "Failed to parse checkInHour", err.Error())
		}

		var socialMediasMap map[string]string
		if len(socialMedias) > 0 {
			socialMediasMap = make(map[string]string)
		}

		for _, sosmed := range socialMedias {
			socialMediasMap[strings.ToLower(sosmed.Platform)] = sosmed.Link
		}

		var name, province, city, subDistrict string
		if strings.TrimSpace(req.Province) != "" {
			province = utils.CapitalizeWords(req.Province)
		}
		if strings.TrimSpace(req.District) != "" {
			city = utils.CapitalizeWords(req.District)
		}
		if strings.TrimSpace(req.SubDistrict) != "" {
			subDistrict = utils.CapitalizeWords(req.SubDistrict)
		}
		if strings.TrimSpace(req.Name) != "" {
			name = utils.CapitalizeWords(req.Name)
		}

		hotel := &entity.Hotel{
			Name:               name,
			AddrSubDistrict:    subDistrict,
			AddrCity:           city,
			AddrProvince:       province,
			IsAPI:              false,
			Description:        req.Description,
			Rating:             req.Rating,
			StatusID:           1,
			CancellationPeriod: hu.config.DefaultCancellationPeriod,
			CheckInHour:        checkInHour,
			CheckOutHour:       checkOutHour,
			SocialMedia:        socialMediasMap,
			Email:              req.Email,
		}

		// Create hotel
		hotel, err = hu.hotelRepo.CreateHotel(txCtx, hotel)
		if err != nil {
			logger.Error(ctx, "Error inserting hotel", err.Error())
			return err
		}

		hotelID = hotel.ID

		// File hotel upload and attachment
		if len(req.Photos) > 0 {
			// Upload photos
			photoURLs, err := hu.uploadMultiple(txCtx, req.Photos, constant.ConstPublic, "hotel", fmt.Sprintf("%d", hotel.ID), "gallery")
			if err != nil {
				logger.Error(ctx, "Error uploading hotel photos", err.Error())
				return err
			}
			hotel.Photos = photoURLs

			// Attach photos to hotel
			if err := hu.hotelRepo.AttachPhotosHotel(txCtx, hotel.ID, photoURLs); err != nil {
				logger.Error(ctx, "Failed to attach hotel photos", err.Error())
				return err
			}

		}

		// Nearby places
		if err := hu.hotelRepo.AttachNearbyPlaces(txCtx, hotel.ID, nearbyPlaces); err != nil {
			logger.Error(ctx, "Failed to attach nearby places", err.Error())
			return err
		}

		// Facilities
		if err := hu.hotelRepo.AttachFacilities(txCtx, hotel.ID, req.Facilities); err != nil {
			logger.Error(ctx, "Failed to attach facilities", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(ctx, "Transaction failed in CreateHotel usecase", err.Error())
		return nil, err
	}

	return &hoteldto.CreateHotelResponse{
		HotelID: hotelID,
	}, nil

}
