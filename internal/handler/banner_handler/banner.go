package banner_handler

import "wtm-backend/internal/domain"

type BannerHandler struct {
	bannerUsecase domain.BannerUsecase
}

func NewBannerHandler(bannerUsecase domain.BannerUsecase) *BannerHandler {
	return &BannerHandler{
		bannerUsecase: bannerUsecase,
	}
}
