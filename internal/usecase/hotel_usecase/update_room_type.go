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
)

func (hu *HotelUsecase) UpdateRoomType(ctx context.Context, req *hoteldto.UpdateRoomTypeRequest) error {
	return hu.dbTransaction.WithTransaction(ctx, func(txCtx context.Context) error {

		var additionalFeatures []hoteldto.RoomAdditional
		if len(req.Additional) > 0 {
			if err := json.Unmarshal([]byte(req.Additional), &additionalFeatures); err != nil {
				logger.Error(txCtx, "Failed to unmarshal AddRoomTypeRequest-additional", err.Error())
				return err
			}
		}

		// Get existing room type
		roomType, err := hu.hotelRepo.GetRoomTypeByID(txCtx, req.RoomTypeID)
		if err != nil {
			logger.Error(txCtx, "Failed to get room type by ID", err.Error())
			return err
		}

		roomType.Name = req.Name
		roomType.IsSmokingAllowed = &req.IsSmokingRoom
		roomType.MaxOccupancy = req.MaxOccupancy
		roomType.RoomSize = req.RoomSize
		roomType.Description = req.Description

		var unchangedAdditions []entity.CustomRoomAdditionalWithID
		for _, id := range req.UnchangedAdditionsIDs {
			for _, addition := range roomType.RoomAdditions {
				if addition.ID == id {
					unchangedAdditions = append(unchangedAdditions, entity.CustomRoomAdditionalWithID{
						ID:         addition.ID,
						Name:       addition.Name,
						Category:   addition.Category,
						Price:      addition.Price,
						Pax:        addition.Pax,
						IsRequired: addition.IsRequired,
					})
					break
				}
			}
		}
		roomType.RoomAdditions = unchangedAdditions

		var fixPhotos []string
		for _, photo := range roomType.Photos {
			for _, roomPhoto := range req.UnchangedRoomPhotos {
				if roomPhoto != "" {
					_, photoUrl, err := hu.fileStorage.ExtractBucketAndObject(txCtx, roomPhoto)
					if err != nil {
						logger.Error(txCtx, "Failed to extract bucket and object from room photo", err.Error())
						continue
					}
					if photoUrl == photo {
						fixPhotos = append(fixPhotos, photo)
						break
					}
				}
			}
		}
		roomType.Photos = fixPhotos

		if len(req.Photos) > 0 {
			photoURLs, err := hu.uploadMultiple(txCtx, req.Photos, constant.ConstPublic, "hotel", fmt.Sprintf("%d", roomType.HotelID), "room_type", req.Name)
			if err != nil {
				logger.Error(txCtx, "Failed to upload room photos", err.Error())
				return err
			}

			roomType.Photos = append(roomType.Photos, photoURLs...)

		}

		var withoutBreakfast hoteldto.BreakfastBase
		if strings.TrimSpace(req.WithoutBreakfast) != "" {
			if err := json.Unmarshal([]byte(req.WithoutBreakfast), &withoutBreakfast); err != nil {
				logger.Error(txCtx, "Failed to unmarshal AddRoomTypeRequest-without_breakfast", err.Error())
				return err
			}

			withoutBreakfastEntity := entity.CustomBreakfastWithID{
				ID:     roomType.WithoutBreakfast.ID,
				Price:  withoutBreakfast.Price,
				IsShow: withoutBreakfast.IsShow,
			}
			roomType.WithoutBreakfast = withoutBreakfastEntity
		}

		var withBreakfast hoteldto.BreakfastWith
		if strings.TrimSpace(req.WithBreakfast) != "" {
			if err := json.Unmarshal([]byte(req.WithBreakfast), &withBreakfast); err != nil {
				logger.Error(txCtx, "Failed to unmarshal AddRoomTypeRequest-with_breakfast", err.Error())
				return err
			}

			withBreakfastEntity := entity.CustomBreakfastWithID{
				ID:     roomType.WithBreakfast.ID,
				Price:  withBreakfast.Price,
				Pax:    withBreakfast.Pax,
				IsShow: withBreakfast.IsShow,
			}
			roomType.WithBreakfast = withBreakfastEntity
		}

		if err := hu.hotelRepo.UpdateRoomType(txCtx, roomType); err != nil {
			logger.Error(txCtx, "Failed to update room type", err.Error())
			return err
		}

		if err := hu.hotelRepo.AttachBedTypesToRoomType(txCtx, roomType.ID, req.BedTypes); err != nil {
			logger.Error(txCtx, "Failed to attach bed types", err.Error())
			return err
		}

		var additionalFeaturesEntity []entity.CustomRoomAdditional
		for _, additional := range additionalFeatures {
			additionalFeaturesEntity = append(additionalFeaturesEntity, entity.CustomRoomAdditional{
				Name:       additional.Name,
				Category:   additional.Category,
				Price:      additional.Price,
				Pax:        additional.Pax,
				IsRequired: additional.IsRequired,
			})
		}

		if err := hu.hotelRepo.AttachRoomAdditions(txCtx, roomType.ID, additionalFeaturesEntity); err != nil {
			logger.Error(txCtx, "Failed to attach facilities", err.Error())
			return err
		}

		return nil
	})
}
