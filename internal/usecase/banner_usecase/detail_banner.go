package banner_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) DetailBanner(ctx context.Context, req *bannerdto.DetailBannerRequest) (*bannerdto.DetailBannerResponse, error) {
	banner, err := bu.bannerRepo.GetBannerByExternalID(ctx, req.BannerID)
	if err != nil {
		logger.Error(ctx, "Error fetching banner details:", err.Error())
		return nil, err
	}

	bucketName := fmt.Sprintf("%s-%s", constant.ConstBanner, constant.ConstPublic)
	bannerUrl, err := bu.fileStorage.GetFile(ctx, bucketName, banner.ImageURL)
	if err != nil {
		logger.Error(ctx, "Error getting banner image", err.Error())

	}

	resp := &bannerdto.DetailBannerResponse{
		Banner: bannerdto.BannerData{
			ID:          req.BannerID,
			Title:       banner.Title,
			Description: banner.Description,
			ImageURL:    bannerUrl,
			IsActive:    banner.IsActive,
		},
	}

	return resp, nil
}
