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

func (hu *HotelUsecase) AddRoomType(ctx context.Context, hotelID uint, req *hoteldto.AddRoomTypeRequest) error {
	return hu.dbTransaction.WithTransaction(ctx, func(txCtx context.Context) error {

		var additionalFeatures []hoteldto.RoomAdditional
		if len(req.Additional) > 0 {
			if err := json.Unmarshal([]byte(req.Additional), &additionalFeatures); err != nil {
				logger.Error(txCtx, "Failed to unmarshal AddRoomTypeRequest-additional", err.Error())
				return err
			}
		}

		// Parse OtherPreferences (simple list of names)
		var otherPreferences []string
		if strings.TrimSpace(req.OtherPreferences) != "" {
			if err := json.Unmarshal([]byte(req.OtherPreferences), &otherPreferences); err != nil {
				logger.Error(txCtx, "Failed to unmarshal AddRoomTypeRequest-other_preferences", err.Error())
				return err
			}
		}

		rt := &entity.RoomType{
			HotelID:          hotelID,
			Name:             req.Name,
			IsSmokingAllowed: &req.IsSmokingRoom,
			MaxOccupancy:     req.MaxOccupancy,
			RoomSize:         req.RoomSize,
			Description:      req.Description,
			TotalUnit:        1,
		}

		rt, err := hu.hotelRepo.CreateRoomType(txCtx, rt)
		if err != nil {
			logger.Error(ctx, "Failed to create room type", err.Error())
			return err
		}

		if len(req.Photos) > 0 {
			photoURLs, err := hu.uploadMultiple(txCtx, req.Photos, constant.ConstPublic, "hotel", fmt.Sprintf("%d", hotelID), "room_type", req.Name)
			if err != nil {
				logger.Error(ctx, "Failed to upload room photos", err.Error())
				return err
			}

			rt.Photos = photoURLs

			if err := hu.hotelRepo.AttachPhotosRoomType(txCtx, rt.ID, photoURLs); err != nil {
				logger.Error(ctx, "Failed to attach room type photos", err.Error())
				return err
			}

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

		if len(additionalFeaturesEntity) > 0 {
			if err := hu.hotelRepo.AttachRoomAdditions(txCtx, rt.ID, additionalFeaturesEntity); err != nil {
				logger.Error(ctx, "Failed to attach facilities", err.Error())
				return err
			}
		}

		// Attach "Other Preferences" if provided
		if len(otherPreferences) > 0 {
			if err := hu.hotelRepo.AttachRoomPreferences(txCtx, rt.ID, otherPreferences); err != nil {
				logger.Error(ctx, "Failed to attach other preferences", err.Error())
				return err
			}
		}

		if len(req.BedTypes) > 0 {
			if err := hu.hotelRepo.AttachBedTypesToRoomType(txCtx, rt.ID, req.BedTypes); err != nil {
				logger.Error(ctx, "Failed to attach bed types", err.Error())
				return err
			}
		}

		var withoutBreakfast hoteldto.BreakfastBase
		if strings.TrimSpace(req.WithoutBreakfast) != "" {
			if err := json.Unmarshal([]byte(req.WithoutBreakfast), &withoutBreakfast); err != nil {
				logger.Error(ctx, "Failed to unmarshal AddRoomTypeRequest-without_breakfast", err.Error())
				return err
			}

			withoutBreakfastEntity := &entity.CustomBreakfast{
				Price:  withoutBreakfast.Price, // DEPRECATED: Keep for backward compatibility
				Prices: withoutBreakfast.Prices,
				IsShow: withoutBreakfast.IsShow,
			}
			// Fallback: if Prices is empty but Price is set, convert Price to Prices
			if len(withoutBreakfastEntity.Prices) == 0 && withoutBreakfastEntity.Price > 0 {
				withoutBreakfastEntity.Prices = map[string]float64{"IDR": withoutBreakfastEntity.Price}
			}

			if err := hu.hotelRepo.CreateRoomPrice(txCtx, rt.ID, withoutBreakfastEntity, false); err != nil {
				logger.Error(ctx, "Failed to create price without breakfast", err.Error())
				return err
			}
		}

		var withBreakfast hoteldto.BreakfastWith
		if strings.TrimSpace(req.WithBreakfast) != "" {
			if err := json.Unmarshal([]byte(req.WithBreakfast), &withBreakfast); err != nil {
				logger.Error(ctx, "Failed to unmarshal AddRoomTypeRequest-with_breakfast", err.Error())
				return err
			}

			withBreakfastEntity := &entity.CustomBreakfast{
				Price:  withBreakfast.Price, // DEPRECATED: Keep for backward compatibility
				Prices: withBreakfast.Prices,
				Pax:    withBreakfast.Pax,
				IsShow: withBreakfast.IsShow,
			}
			// Fallback: if Prices is empty but Price is set, convert Price to Prices
			if len(withBreakfastEntity.Prices) == 0 && withBreakfastEntity.Price > 0 {
				withBreakfastEntity.Prices = map[string]float64{"IDR": withBreakfastEntity.Price}
			}

			if err := hu.hotelRepo.CreateRoomPrice(txCtx, rt.ID, withBreakfastEntity, true); err != nil {
				logger.Error(ctx, "Failed to create price with breakfast", err.Error())
				return err
			}
		}

		return nil
	})
}
