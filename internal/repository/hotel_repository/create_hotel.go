package hotel_repository

import (
	"context"
	"gorm.io/datatypes"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (hr *HotelRepository) CreateHotel(ctx context.Context, hotel *entity.Hotel) (*entity.Hotel, error) {
	db := hr.db.GetTx(ctx)

	var hotelModel model.Hotel
	if err := utils.CopyStrict(&hotelModel, hotel); err != nil {
		logger.Error(ctx, "Failed to copy hotel entity to model", err.Error())
		return nil, err
	}

	if hotel.SocialMedia != nil {
		jsonBytes, err := utils.MapToJSON(hotel.SocialMedia)
		if err != nil {
			logger.Error(ctx, "Failed to map hotel social media", err.Error())
		} else {
			hotelModel.SocialMedia = datatypes.JSON(jsonBytes)
		}
	}

	if err := db.WithContext(ctx).Create(&hotelModel).Error; err != nil {
		logger.Error(ctx, "Failed to create hotel", err.Error())
		return nil, err
	}

	// Copy back the result including Id
	if err := utils.CopyStrict(hotel, &hotelModel); err != nil {
		logger.Error(ctx, "Failed to copy created hotel back to entity", err.Error())
		return nil, err
	}

	return hotel, nil
}
