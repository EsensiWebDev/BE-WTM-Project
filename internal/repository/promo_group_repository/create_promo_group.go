package promo_group_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) CreatePromoGroup(ctx context.Context, promoGroup *entity.PromoGroup) error {
	db := pgr.db.GetTx(ctx)

	var promoGroupModel model.PromoGroup
	if err := utils.CopyStrict(&promoGroupModel, promoGroup); err != nil {
		logger.Error(ctx, "Error copying promo group entity to model", err.Error())
		return err
	}

	err := db.WithContext(ctx).
		Create(&promoGroupModel).Error
	if err != nil {
		return err
	}
	return nil
}
