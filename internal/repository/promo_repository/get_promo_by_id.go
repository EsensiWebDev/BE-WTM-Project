package promo_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) GetPromoByID(ctx context.Context, promoID uint, selectedFields []string) (*entity.Promo, error) {
	db := pr.db.GetTx(ctx)

	var promo model.Promo
	query := db.WithContext(ctx).Model(&model.Promo{})

	if len(selectedFields) > 0 {
		query = query.Select(selectedFields)
	} else {
		query = query.Preload("PromoType").
			Preload("PromoGroups").
			Preload("PromoRoomTypes").
			Preload("PromoRoomTypes.RoomType").
			Preload("PromoRoomTypes.RoomType.Hotel")
	}

	if err := query.Where("id = ?", promoID).First(&promo).Error; err != nil {
		if pr.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "Promo not found with Id", promoID)
			return nil, nil
		}
		logger.Error(ctx, "Error finding promo by Id", err.Error())
		return nil, err
	}

	var promoEntity entity.Promo
	if err := utils.CopyStrict(&promoEntity, &promo); err != nil {
		logger.Error(ctx, "Error copying promo model to entity", err.Error())
		return nil, err
	}

	var detailPromo entity.PromoDetail
	if err := json.Unmarshal(promo.Detail, &detailPromo); err != nil {
		logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
	}
	promoEntity.Detail = detailPromo
	promoEntity.ExternalID = promo.ExternalID.ExternalID

	if len(selectedFields) == 0 {
		promoEntity.PromoTypeName = promo.PromoType.Name

		for i, promoRoomTypeE := range promoEntity.PromoRoomTypes {
			for _, promoRoomTypesM := range promo.PromoRoomTypes {
				if promoRoomTypeE.RoomTypeID == promoRoomTypesM.RoomTypeID {
					promoRoomTypeE.RoomTypeName = promoRoomTypesM.RoomType.Name
					promoRoomTypeE.HotelID = promoRoomTypesM.RoomType.HotelID
					promoRoomTypeE.HotelName = promoRoomTypesM.RoomType.Hotel.Name
				}
			}
			promoEntity.PromoRoomTypes[i] = promoRoomTypeE
		}
	}

	return &promoEntity, nil
}
