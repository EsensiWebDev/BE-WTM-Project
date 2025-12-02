package banner_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) ListActiveBanners(ctx context.Context) (*bannerdto.ListActiveBannerResponse, error) {
	filterRepo := &filter.BannerFilter{}
	filterRepo.IsActive = new(bool)
	*filterRepo.IsActive = true

	banners, _, err := bu.bannerRepo.GetBanners(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Failed to get active banners", err.Error())
		return nil, err
	}

	var activeBanner []bannerdto.ActiveBanner
	activeBanner = make([]bannerdto.ActiveBanner, 0, len(banners))
	for _, banner := range banners {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstBanner, constant.ConstPublic)
		bannerUrl, err := bu.fileStorage.GetFile(ctx, bucketName, banner.ImageURL)
		if err != nil {
			logger.Error(ctx, "Error getting banner image", err.Error())
		}
		activeBanner = append(activeBanner, bannerdto.ActiveBanner{
			ID:          banner.ExternalID,
			Title:       banner.Title,
			Description: banner.Description,
			ImageURL:    bannerUrl,
		})
	}

	return &bannerdto.ListActiveBannerResponse{
		Banners: activeBanner,
	}, nil
}
