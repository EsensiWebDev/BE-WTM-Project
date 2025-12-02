package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/repository/filter"
)

type BannerUsecase interface {
	UpsertBanner(ctx context.Context, req *bannerdto.UpsertBannerRequest, reqID *bannerdto.DetailBannerRequest) error
	ListBanners(ctx context.Context, req *bannerdto.ListBannerRequest) (*bannerdto.ListBannerResponse, error)
	ListActiveBanners(ctx context.Context) (*bannerdto.ListActiveBannerResponse, error)
	DetailBanner(ctx context.Context, req *bannerdto.DetailBannerRequest) (*bannerdto.DetailBannerResponse, error)
	RemoveBanner(ctx context.Context, req *bannerdto.DetailBannerRequest) error
	UpdateStatusBanner(ctx context.Context, req *bannerdto.UpdateStatusBannerRequest) error
	UpdateOrderBanner(ctx context.Context, req *bannerdto.UpdateOrderBannerRequest) error
}

type BannerRepository interface {
	CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	GetBannerByID(ctx context.Context, id uint) (*entity.Banner, error)
	GetBanners(ctx context.Context, filter *filter.BannerFilter) ([]entity.Banner, int64, error)
	DeleteBanner(ctx context.Context, id uint) error
	UpdateStatusBanner(ctx context.Context, id string, isActive bool) error
	UpdateOrderBanner(ctx context.Context, id string, dir string) error
	GetBannerByExternalID(ctx context.Context, externalID string) (*entity.Banner, error)
}
