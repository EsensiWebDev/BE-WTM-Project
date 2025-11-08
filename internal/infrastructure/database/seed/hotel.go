package seed

import (
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
)

func (s *Seed) SeedHotel() {

	var countStatusHotel int64
	s.db.Model(&model.StatusHotel{}).Count(&countStatusHotel)
	if countStatusHotel == 0 {
		statusHotels := []model.StatusHotel{
			{Status: constant.StatusHotelInReview},
			{Status: constant.StatusHotelApproved},
			{Status: constant.StatusHotelRejected},
		}

		if err := s.db.Create(&statusHotels).Error; err != nil {
			log.Fatalf("Failed to seed status hotels: %s", err.Error())
		}
		log.Println("Seeding status hotels completed")
	}
}
