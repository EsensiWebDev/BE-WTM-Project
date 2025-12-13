package promo_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pr *PromoRepository) GetPromoByExternalID(ctx context.Context, externalID string) (*entity.Promo, error) {
	db := pr.db.GetTx(ctx)

	var promo model.Promo
	err := db.Where("external_id = ?", externalID).First(&promo).Error
	if err != nil {
		logger.Error(ctx, "GetPromoByExternalID", err.Error())
		return nil, err
	}

	var entityPromo entity.Promo
	if err := utils.CopyStrict(&entityPromo, &promo); err != nil {
		logger.Error(ctx, "GetPromoByExternalID", err.Error())
		return nil, err
	}

	return &entityPromo, nil
}
