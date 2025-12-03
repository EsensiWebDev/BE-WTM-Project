package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/repository/filter"
)

type PromoGroupUsecase interface {
	CreatePromoGroup(ctx context.Context, name string) error
	DetailPromoGroup(ctx context.Context, promoGroupID uint) (*entity.PromoGroup, error)
	ListPromoGroups(ctx context.Context, req *promogroupdto.ListPromoGroupRequest) (*promogroupdto.ListPromoGroupResponse, int64, error)
	ListPromoGroupPromos(ctx context.Context, req *promogroupdto.ListPromoGroupPromosRequest) (*promogroupdto.ListPromoGroupPromosResponse, int64, error)
	ListPromoGroupMembers(ctx context.Context, req *promogroupdto.ListPromoGroupMemberRequest) (*promogroupdto.ListPromoGroupMemberResponse, int64, error)
	AssignPromoGroupMember(ctx context.Context, req *promogroupdto.AssignPromoGroupMemberRequest) error
	RemovePromoGroupMember(ctx context.Context, req *promogroupdto.RemovePromoGroupMemberRequest) error
	AssignPromoToGroup(ctx context.Context, req *promogroupdto.AssignPromoToGroupRequest) error
	RemovePromoFromGroup(ctx context.Context, req *promogroupdto.RemovePromoFromGroupRequest) error
	ListUnassignedPromos(ctx context.Context, req *promogroupdto.ListUnassignedPromosRequest) (*promogroupdto.ListUnassignedPromosResponse, error)
	RemovePromoGroup(ctx context.Context, promoGroupID uint) error
}

type PromoGroupRepository interface {
	GetPromoGroupByID(ctx context.Context, promoGroupID uint) (*entity.PromoGroup, error)
	GetPromoGroups(ctx context.Context, search string, limit, page int) ([]entity.PromoGroup, int64, error)
	CreatePromoGroup(ctx context.Context, promoGroup *entity.PromoGroup) error
	GetPromoGroupMembers(ctx context.Context, promoGroupID uint, limit, page int) ([]entity.User, int64, error)
	RemovePromoGroupMember(ctx context.Context, promoGroupID uint, memberID uint) error
	GetPromosByPromoGroupID(ctx context.Context, promoGroupID uint, search string, limit, page int) ([]entity.Promo, int64, error)
	AssignPromoToGroup(ctx context.Context, promoGroupID uint, promoID uint) error
	RemovePromoFromGroup(ctx context.Context, promoGroupID uint, promoID uint) error
	GetUnassignedPromos(ctx context.Context, filterReq *filter.PromoGroupFilter) ([]entity.Promo, int64, error)
	DeletePromoGroup(ctx context.Context, promoGroupID uint) error
	CheckPromoGroupExists(ctx context.Context, name string) bool
}
