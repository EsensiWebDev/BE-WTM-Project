package seed

import (
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
)

func (s *Seed) SeedBooking() {

	var countStatusBooking, countStatusPayment int64
	s.db.Model(&model.StatusBooking{}).Count(&countStatusBooking)
	s.db.Model(&model.StatusPayment{}).Count(&countStatusPayment)
	if countStatusBooking == 0 {
		statusBookings := []model.StatusBooking{
			{Status: constant.StatusBookingInCart},
			{Status: constant.StatusBookingInReview},
			{Status: constant.StatusBookingApproved},
			{Status: constant.StatusBookingRejected},
		}

		if err := s.db.Create(&statusBookings).Error; err != nil {
			log.Fatalf("Failed to seed status bookings: %s", err.Error())
		}
		log.Println("Seeding status booking completed")
	}
	if countStatusPayment == 0 {
		statusPayments := []model.StatusPayment{
			{Status: constant.StatusPaymentUnpaid},
			{Status: constant.StatusPaymentPaid},
		}

		if err := s.db.Create(&statusPayments).Error; err != nil {
			log.Fatalf("Failed to seed status payments: %s", err.Error())
		}
		log.Println("Seeding status payments completed")
	}
}
