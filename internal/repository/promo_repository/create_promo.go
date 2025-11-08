package promo_repository

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) CreatePromo(ctx context.Context, promo *entity.Promo) error {
	db := pr.db.GetTx(ctx)

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

	// Inject promo groups association secara langsung
	if len(promo.PromoGroupIDs) > 0 {
		for _, id := range promo.PromoGroupIDs {
			promoModel.PromoGroups = append(promoModel.PromoGroups, model.PromoGroup{
				Model: gorm.Model{ID: id},
			})
		}
	}

	// Create promo + auto insert PromoRoomTypes + PromoGroups via GORM associations
	if err := db.Create(&promoModel).Error; err != nil {
		logger.Error(ctx, "Error creating promo with associations", err.Error())
		return err
	}

	return nil
}
