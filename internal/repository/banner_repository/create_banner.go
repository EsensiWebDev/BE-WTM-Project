package banner_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BannerRepository) CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error) {
	db := br.db.GetTx(ctx)

	var bannerModel model.Banner
	if err := utils.CopyStrict(&bannerModel, &banner); err != nil {
		logger.Error(ctx, "Error copying banner", err.Error())
		return nil, err
	}

	if err := db.Create(&bannerModel).Error; err != nil {
		logger.Error(ctx, "Error creating banner in database", err.Error())
		return nil, err
	}

	if err := utils.CopyStrict(banner, &bannerModel); err != nil {
		logger.Error(ctx, "Error copying banner model to entity", err.Error())
		return nil, err
	}

	return banner, nil
}
