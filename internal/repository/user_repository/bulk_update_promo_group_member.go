package user_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ur *UserRepository) BulkUpdatePromoGroupMember(ctx context.Context, userIDs []uint, promoGroupID uint) error {
	db := ur.db.GetTx(ctx)

	if len(userIDs) == 0 {
		logger.Warn(ctx,
			"No user IDs provided for bulk update")
		return nil
	}

	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id IN (?)", userIDs).
		Update("promo_group_id", promoGroupID).Error

	if err != nil {
		logger.Error(ctx, "Error bulk updating promo group", err.Error())
		return err
	}

	return nil
}
