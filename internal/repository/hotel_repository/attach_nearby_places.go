package hotel_repository

import (
	"context"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (hr *HotelRepository) AttachNearbyPlaces(ctx context.Context, hotelID uint, nearbyPlaces []hoteldto.NearbyPlace) error {
	db := hr.db.GetTx(ctx)

	for _, np := range nearbyPlaces {
		placeModel := model.NearbyPlace{
			Name: np.Name,
		}

		// Ensure NearbyPlace exists or create it
		if err := db.WithContext(ctx).
			Where("name = ?", placeModel.Name).
			FirstOrCreate(&placeModel).Error; err != nil {
			logger.Error(ctx, "Failed to attach nearby places", err.Error())
			return err
		}

		// Link hotel and place
		link := model.HotelNearbyPlace{
			HotelID:       hotelID,
			NearbyPlaceID: placeModel.ID,
			Radius:        np.Distance,
		}

		if err := db.WithContext(ctx).Create(&link).Error; err != nil {
			logger.Error(ctx, "Failed to attach nearby place to hotel", err.Error())
			return err
		}
	}

	return nil
}
