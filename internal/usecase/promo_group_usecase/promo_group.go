package promo_group_usecase

import "wtm-backend/internal/domain"

type PromoGroupUsecase struct {
	promoGroupRepo domain.PromoGroupRepository
	userRepo       domain.UserRepository
}

func NewPromoGroupUsecase(promoGroupRepo domain.PromoGroupRepository, userRepo domain.UserRepository) *PromoGroupUsecase {
	return &PromoGroupUsecase{
		promoGroupRepo: promoGroupRepo,
		userRepo:       userRepo,
	}
}
