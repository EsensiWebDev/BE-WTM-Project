package promo_group_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) GetPromoGroupByID(ctx context.Context, promoGroupID uint) (*entity.PromoGroup, error) {
	db := pgr.db.GetTx(ctx)

	var promoGroup model.PromoGroup
	err := db.WithContext(ctx).
		Select("id,name").
		Where("id = ?", promoGroupID).
		First(&promoGroup).Error
	if err != nil {
		if pgr.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "Promo group not found with Id", promoGroupID)
			return nil, nil
		}
		logger.Error(ctx, "Error finding promo group by Id", err.Error())
		return nil, err
	}

	var promoGroupEntity entity.PromoGroup
	if err := utils.CopyPatch(&promoGroupEntity, &promoGroup); err != nil {
		logger.Error(ctx, "Error copying promo group model to entity", err.Error())
		return nil, err
	}

	return &promoGroupEntity, nil
}
