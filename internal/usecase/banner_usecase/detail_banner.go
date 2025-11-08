package banner_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) DetailBanner(ctx context.Context, id uint) (*entity.Banner, error) {
	banner, err := bu.bannerRepo.GetBannerByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "Error fetching banner details:", err.Error())
		return nil, err
	}
	return banner, nil
}
