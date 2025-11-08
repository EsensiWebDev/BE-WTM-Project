package banner_usecase

import "wtm-backend/internal/domain"

type BannerUsecase struct {
	bannerRepo  domain.BannerRepository
	dbTrx       domain.DatabaseTransaction
	fileStorage domain.StorageClient
}

func NewBannerUsecase(bannerRepo domain.BannerRepository, dbTrx domain.DatabaseTransaction, fileStorage domain.StorageClient) *BannerUsecase {
	return &BannerUsecase{
		bannerRepo:  bannerRepo,
		dbTrx:       dbTrx,
		fileStorage: fileStorage,
	}
}
