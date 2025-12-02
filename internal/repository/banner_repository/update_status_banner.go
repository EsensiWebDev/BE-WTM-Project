package banner_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BannerRepository) UpdateStatusBanner(ctx context.Context, id string, isActive bool) error {
	db := br.db.GetTx(ctx)

	err := db.WithContext(ctx).
		Model(&model.Banner{}).
		Where("external_id = ?", id).
		Update("is_active", isActive).Error
	if err != nil {
		logger.Error(ctx, "Error updating banner status", err.Error())
		return err
	}

	return nil
}
