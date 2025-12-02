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
		if req.NearbyPlaces != "" {
			if err := json.Unmarshal([]byte(req.NearbyPlaces), &nearbyPlaces); err != nil {
				logger.Error(ctx, "Failed to unmarshal UpdateHotel-NearbyPlaces", err.Error())
				return err
			}
		}

		var socialMedias []hoteldto.SocialMedia
		if req.SocialMedias != "" {
			if err := json.Unmarshal([]byte(req.SocialMedias), &socialMedias); err != nil {
				logger.Error(ctx, "Failed to unmarshal UpdateHotel-SocialMedias", err.Error())
				return err
			}
		}
		hotel, err := hu.hotelRepo.GetHotelByID(txCtx, req.HotelID, constant.RoleAdmin)
		if err != nil {
			logger.Error(ctx, "failed to get hotel by ID", err.Error())
			return fmt.Errorf("failed to get hotel by ID: %w", err)
		}

		var photoHotel []string
		for _, photoOri := range hotel.Photos {
			for _, photo := range req.UnchangedHotelPhotos {
				if photo != "" {
					_, photoURL, err := hu.fileStorage.ExtractBucketAndObject(txCtx, photo)
					if err != nil {
						logger.Error(ctx, "failed to extract bucket and object from unchanged hotel photo", err.Error())
						continue
					}
					if photoURL == photoOri {
						photoHotel = append(photoHotel, photoURL)
						break
					}
				}
			}
		}
		hotel.Photos = photoHotel

		hotel.Name = req.Name
		hotel.AddrSubDistrict = req.SubDistrict
		hotel.AddrCity = req.District
		hotel.AddrProvince = req.Province
		hotel.Description = req.Description
		hotel.Rating = req.Rating
		hotel.Email = req.Email

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
		if len(req.Photos) > 0 {
			// Upload photos
			photoURLs, err := hu.uploadMultiple(txCtx, req.Photos, constant.ConstPublic, "hotel", fmt.Sprintf("%d", hotel.ID), "gallery")
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
		if err := hu.hotelRepo.AttachFacilities(txCtx, hotel.ID, req.Facilities); err != nil {
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
