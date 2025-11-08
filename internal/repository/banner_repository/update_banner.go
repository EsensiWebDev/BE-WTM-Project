package banner_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BannerRepository) UpdateBanner(ctx context.Context, banner *entity.Banner) error {
	db := br.db.GetTx(ctx)

	var bannerModel model.Banner
	if err := utils.CopyStrict(&bannerModel, &banner); err != nil {
		logger.Error(ctx, "Error copying banner", err.Error())
		return err
	}

	if err := db.Model(&bannerModel).Where("id = ?", banner.ID).Updates(bannerModel).Error; err != nil {
		logger.Error(ctx, "Error updating banner in database", err.Error())
		return err
	}

	return nil
}
