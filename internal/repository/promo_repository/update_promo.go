package promo_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) UpdatePromo(ctx context.Context, promo *entity.Promo) error {
	db := pr.db.GetTx(ctx)

	// Step 1: Build promoModel dari entity
	var promoModel model.Promo
	if err := utils.CopyStrict(&promoModel, promo); err != nil {
		logger.Error(ctx, "Error copying promo entity to model", err.Error())
		return err
	}

	jsonDetail, err := json.Marshal(promo.Detail)
	if err != nil {
		logger.Error(ctx, "Error marshalling promo detail to JSON", err.Error())
	}
	if jsonDetail != nil {
		promoModel.Detail = jsonDetail
	}

	// Step 2: Update promo utama
	if err := db.Model(&model.Promo{}).
		Where("id = ?", promo.ID).
		Updates(&promoModel).Error; err != nil {
		logger.Error(ctx, "Error updating promo", err.Error())
		return err
	}

	// Step 3: Update PromoRoomTypes → Delete & Insert ulang
	if err := db.Where("promo_id = ?", promo.ID).
		Delete(&model.PromoRoomType{}).Error; err != nil {
		logger.Error(ctx, "Error deleting promo room types", err.Error())
		return err
	}

	if len(promo.PromoRoomTypes) > 0 {
		var roomTypeModels []model.PromoRoomType
		for _, rt := range promo.PromoRoomTypes {
			roomTypeModels = append(roomTypeModels, model.PromoRoomType{
				PromoID:     promo.ID,
				RoomTypeID:  rt.RoomTypeID,
				TotalNights: rt.TotalNights,
			})
		}
		if err := db.Create(&roomTypeModels).Error; err != nil {
			logger.Error(ctx, "Error inserting promo room types", err.Error())
			return err
		}
	}

	return nil
}
