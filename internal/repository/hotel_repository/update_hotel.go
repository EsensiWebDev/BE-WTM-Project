package hotel_repository

import (
	"context"
	"encoding/json"
	"github.com/lib/pq"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) UpdateHotel(ctx context.Context, hotel *entity.Hotel) error {
	db := hr.db.GetTx(ctx)

	var hotelModel model.Hotel
	if err := utils.CopyPatch(&hotelModel, hotel); err != nil {
		logger.Error(ctx, "failed to copy hotel", err.Error())
		return err
	}

	hotelModel.RoomTypes = nil
	hotelModel.Facilities = nil
	hotelModel.HotelNearbyPlaces = nil

	if hotel.IsAPI {
		hotelModel.IsAPI = true
	}

	// Manual conversion for incompatible fields
	hotelModel.Photos = pq.StringArray(hotel.Photos)

	if hotel.SocialMedia != nil {
		socialJSON, err := json.Marshal(hotel.SocialMedia)
		if err != nil {
			logger.Error(ctx, "failed to marshal social media", err.Error())
			return err
		}
		hotelModel.SocialMedia = socialJSON
	}

	// Proceed with update
	if err := db.WithContext(ctx).Model(&model.Hotel{}).
		Where("id = ?", hotel.ID).
		Updates(&hotelModel).Error; err != nil {
		logger.Error(ctx, "failed to update hotel", err.Error())
		return err
	}

	// Hapus semua relasi lama
	if err := db.WithContext(ctx).
		Unscoped().
		Where("hotel_id = ?", hotel.ID).Debug().
		Delete(&model.HotelNearbyPlace{}).Error; err != nil {
		logger.Error(ctx, "failed to clear hotel nearby places", err.Error())
		return err
	}

	// Tambahkan relasi baru dengan radius
	for _, np := range hotel.NearbyPlaces {
		link := model.HotelNearbyPlace{
			HotelID:       hotel.ID,
			NearbyPlaceID: np.ID,
			Radius:        np.Radius,
		}
		if err := db.WithContext(ctx).Debug().Create(&link).Error; err != nil {
			logger.Error(ctx, "failed to attach nearby place to hotel", err.Error())
			return err
		}
	}

	return nil

}
