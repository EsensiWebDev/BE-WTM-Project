package seed

import (
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
)

func (s *Seed) SeedPromo() {

	var countPromoType int64
	s.db.Model(&model.PromoType{}).Count(&countPromoType)
	if countPromoType == 0 {
		promoTypes := []model.PromoType{
			{Name: constant.PromoTypeDiscount},
			{Name: constant.PromoTypeFixedPrice},
			{Name: constant.PromoTypeRoomUpgrade},
			{Name: constant.PromoTypeBenefit},
		}

		if err := s.db.Create(&promoTypes).Error; err != nil {
			log.Fatalf("Failed to seed promo types: %s", err.Error())
		}
		log.Println("Seeding promo types completed")
	}
}
