package banner_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BannerRepository) UpdateOrderBanner(ctx context.Context, id uint, order int) error {
	db := br.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Model(&model.Banner{}).
		Where("id = ?", id).
		Update("display_order", order).Error
	if err != nil {
		logger.Error(ctx, "Error updating banner order", err.Error())
		return err
	}

	return nil
}
