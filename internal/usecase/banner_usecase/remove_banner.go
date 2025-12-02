package banner_usecase

import (
	"context"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) RemoveBanner(ctx context.Context, req *bannerdto.DetailBannerRequest) error {
	banner, err := bu.bannerRepo.GetBannerByExternalID(ctx, req.BannerID)
	if err != nil {
		logger.Error(ctx, "Error fetching banner details:", err.Error())
		return err
	}

	if err = bu.bannerRepo.DeleteBanner(ctx, banner.ID); err != nil {
		logger.Error(ctx, "Error deleting banner:", err.Error())
		return err
	}

	return nil
}
