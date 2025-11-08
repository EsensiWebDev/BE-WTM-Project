package banner_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BannerUsecase) ListBanners(ctx context.Context, req *bannerdto.ListBannerRequest) (*bannerdto.ListBannerResponse, error) {

	filterRepo := &filter.BannerFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest
	filterRepo.IsActive = req.IsActive

	banners, total, err := bu.bannerRepo.GetBanners(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting banners", err.Error())
		return nil, err
	}

	for i := range banners {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstBanner, constant.ConstPublic)
		bannerUrl, err := bu.fileStorage.GetFile(ctx, bucketName, banners[i].ImageURL)
		if err != nil {
			logger.Error(ctx, "Error getting banner image", err.Error())
			continue
		}
		banners[i].ImageURL = bannerUrl
	}

	response := &bannerdto.ListBannerResponse{
		Banners: banners,
		Total:   total,
	}

	return response, nil
}
