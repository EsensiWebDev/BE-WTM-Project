package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/repository/filter"
)

type BannerUsecase interface {
	UpsertBanner(ctx context.Context, req *bannerdto.UpsertBannerRequest, bannerId *uint) error
	ListBanners(ctx context.Context, req *bannerdto.ListBannerRequest) (*bannerdto.ListBannerResponse, error)
	DetailBanner(ctx context.Context, id uint) (*entity.Banner, error)
	RemoveBanner(ctx context.Context, id uint) error
	UpdateStatusBanner(ctx context.Context, req *bannerdto.UpdateStatusBannerRequest) error
	UpdateOrderBanner(ctx context.Context, req *bannerdto.UpdateOrderBannerRequest) error
}

type BannerRepository interface {
	CreateBanner(ctx context.Context, banner *entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, banner *entity.Banner) error
	GetBannerByID(ctx context.Context, id uint) (*entity.Banner, error)
	GetBanners(ctx context.Context, filter *filter.BannerFilter) ([]entity.Banner, int64, error)
	DeleteBanner(ctx context.Context, id uint) error
	UpdateStatusBanner(ctx context.Context, id uint, isActive bool) error
	UpdateOrderBanner(ctx context.Context, id uint, order int) error
}
