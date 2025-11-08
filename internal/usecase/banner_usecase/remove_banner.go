package banner_usecase

import (
	"context"
)

func (bu *BannerUsecase) RemoveBanner(ctx context.Context, id uint) error {
	err := bu.bannerRepo.DeleteBanner(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
