package banner_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BannerRepository) DeleteBanner(ctx context.Context, id uint) error {
	db := br.db.GetTx(ctx)

	// Delete the banner by Id
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.Banner{}).Error; err != nil {
		logger.Error(ctx, "Error deleting banner", err.Error())
		return err
	}

	return nil
}
