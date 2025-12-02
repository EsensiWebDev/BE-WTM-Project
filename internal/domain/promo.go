package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/repository/filter"
)

type PromoUsecase interface {
	ListPromoTypes(ctx context.Context, req *promodto.ListPromoTypesRequest) (*promodto.ListPromoTypesResponse, int64, error)
	UpsertPromo(ctx context.Context, req *promodto.UpsertPromoRequest, promoID *uint) error
	ListPromos(ctx context.Context, req *promodto.ListPromosRequest) (*promodto.ListPromosResponse, error)
	ListPromosForAgent(ctx context.Context, req *promodto.ListPromosForAgentRequest) (*promodto.ListPromosForAgentResponse, error)
	SetStatusPromo(ctx context.Context, req *promodto.SetStatusPromoRequest) error
	PromoByID(ctx context.Context, promoID uint) (*entity.Promo, error)
	RemovePromo(ctx context.Context, promoID uint) error
}

type PromoRepository interface {
	GetPromoTypes(ctx context.Context, filter *filter.DefaultFilter) ([]entity.PromoType, int64, error)
	CreatePromo(ctx context.Context, promo *entity.Promo) error
	GetPromos(ctx context.Context, filterReq *filter.DefaultFilter) ([]entity.Promo, int64, error)
	GetPromosWithHotels(ctx context.Context, filterReq *filter.PromoFilter) ([]entity.Promo, int64, error)
	GetPromoByID(ctx context.Context, promoID uint, selectedFields []string) (*entity.Promo, error)
	UpdatePromoStatus(ctx context.Context, promoID uint, isActive bool) error
	DeletePromo(ctx context.Context, promoID uint) error
	UpdatePromo(ctx context.Context, promo *entity.Promo) error
}
