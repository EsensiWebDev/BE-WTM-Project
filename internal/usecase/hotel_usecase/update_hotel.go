package hotel_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (hu *HotelUsecase) UpdateHotel(ctx context.Context, req *hoteldto.UpdateHotelRequest) error {
	return hu.dbTransaction.WithTransaction(ctx, func(txCtx context.Context) error {
		var nearbyPlaces []hoteldto.NearbyPlace
		if req.CreateHotel.NearbyPlaces != "" {
			if err := json.Unmarshal([]byte(req.CreateHotel.NearbyPlaces), &nearbyPlaces); err != nil {
				logger.Error(ctx, "Failed to unmarshal UpdateHotel-NearbyPlaces", err.Error())
				return err
			}
		}

		var socialMedias []hoteldto.SocialMedia
		if req.CreateHotel.SocialMedias != "" {
			if err := json.Unmarshal([]byte(req.CreateHotel.SocialMedias), &socialMedias); err != nil {
				logger.Error(ctx, "Failed to unmarshal UpdateHotel-SocialMedias", err.Error())
				return err
			}
		}
		hotel, err := hu.hotelRepo.GetHotelByID(txCtx, req.HotelID, constant.RoleAdmin)
		if err != nil {
			logger.Error(ctx, "failed to get hotel by ID", err.Error())
			return fmt.Errorf("failed to get hotel by ID: %w", err)
		}

		hotel.Name = req.CreateHotel.Name
		hotel.AddrSubDistrict = req.CreateHotel.SubDistrict
		hotel.AddrCity = req.CreateHotel.District
		hotel.AddrProvince = req.CreateHotel.Province
		hotel.Description = req.CreateHotel.Description
		hotel.Rating = req.CreateHotel.Rating
		hotel.Email = req.CreateHotel.Email
		hotel.Photos = req.UnchangedHotelPhotos

		socialMediasMap := hotel.SocialMedia
		for _, sosmed := range socialMedias {
			socialMediasMap[sosmed.Platform] = sosmed.Link
		}
		hotel.SocialMedia = socialMediasMap

		var nearbyPlacesEntity []entity.NearbyPlace
		for _, nearbyPlaceID := range req.UnchangedNearbyPlaceIDs {
			for _, place := range hotel.NearbyPlaces {
				if place.ID == nearbyPlaceID {
					nearbyPlacesEntity = append(nearbyPlacesEntity, place)
					break
				}
			}
		}
		hotel.NearbyPlaces = nearbyPlacesEntity

		// File hotel upload and attachment
		if len(req.CreateHotel.Photos) > 0 {
			// Upload photos
			photoURLs, err := hu.uploadMultiple(txCtx, req.CreateHotel.Photos, constant.ConstPublic, "hotel", fmt.Sprintf("%d", hotel.ID), "gallery")
			if err != nil {
				logger.Error(ctx, "Error uploading hotel photos", err.Error())
				return err
			}
			hotel.Photos = append(hotel.Photos, photoURLs...)

		}

		if err := hu.hotelRepo.UpdateHotel(txCtx, hotel); err != nil {
			logger.Error(ctx, "failed to update hotel", err.Error())
			return fmt.Errorf("failed to update hotel: %w", err)
		}

		// Facilities
		if err := hu.hotelRepo.AttachFacilities(txCtx, hotel.ID, req.CreateHotel.Facilities); err != nil {
			logger.Error(ctx, "Failed to attach facilities", err.Error())
			return err
		}

		// Nearby places
		if len(nearbyPlaces) > 0 {
			if err := hu.hotelRepo.AttachNearbyPlaces(txCtx, hotel.ID, nearbyPlaces); err != nil {
				logger.Error(ctx, "failed to attach nearby places", err.Error())
				return fmt.Errorf("failed to attach nearby places: %w", err)
			}
		}

		return nil
	})
}
