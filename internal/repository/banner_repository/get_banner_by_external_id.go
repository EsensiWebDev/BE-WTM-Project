package banner_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BannerRepository) GetBannerByExternalID(ctx context.Context, externalID string) (*entity.Banner, error) {
	db := br.db.GetTx(ctx)

	var banner model.Banner

	if err := db.Where("external_id = ?", externalID).First(&banner).Error; err != nil {
		if br.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "Banner not found")
			return nil, nil // Return nil if no record found
		}
		logger.Error(ctx, "Error getting banner by Id", err.Error())
		return nil, err // Return error if any other error occurs
	}

	var bannerEntity entity.Banner
	if err := utils.CopyStrict(&bannerEntity, banner); err != nil {
		logger.Error(ctx, "Error copying banner model to entity", err.Error())
		return nil, err // Return error if copying fails
	}

	return &bannerEntity, nil
}
