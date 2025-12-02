package seed

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
)

func (s *Seed) SeedBooking() {
	// definisi status booking sesuai kebutuhan client
	statusBookings := []model.StatusBooking{
		{ID: constant.StatusBookingInCartID, Status: constant.StatusBookingInCart},
		{ID: constant.StatusBookingWaitingApprovalID, Status: constant.StatusBookingWaitingApproval},
		{ID: constant.StatusBookingConfirmedID, Status: constant.StatusBookingConfirmed},
		{ID: constant.StatusBookingRejectedID, Status: constant.StatusBookingRejected},
		{ID: constant.StatusBookingCanceledID, Status: constant.StatusBookingCanceled},
	}

	// definisi status payment
	statusPayments := []model.StatusPayment{
		{ID: constant.StatusPaymentUnpaidID, Status: constant.StatusPaymentUnpaid},
		{ID: constant.StatusPaymentPaidID, Status: constant.StatusPaymentPaid},
	}

	// sinkronisasi status booking
	for _, sb := range statusBookings {
		var existing model.StatusBooking
		if err := s.db.First(&existing, sb.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// insert baru
				if err := s.db.Create(&sb).Error; err != nil {
					log.Fatalf("Failed to insert status booking: %s", err.Error())
				}
			} else {
				log.Fatalf("Error checking status booking: %s", err.Error())
			}
		} else {
			// update kalau ada perubahan
			if existing.Status != sb.Status {
				if err := s.db.Model(&existing).Update("status", sb.Status).Error; err != nil {
					log.Fatalf("Failed to update status booking: %s", err.Error())
				}
			}
		}
	}
	log.Println("Seeding status booking completed (sync)")

	// sinkronisasi status payment
	for _, sp := range statusPayments {
		var existing model.StatusPayment
		if err := s.db.First(&existing, sp.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// insert baru
				if err := s.db.Create(&sp).Error; err != nil {
					log.Fatalf("Failed to insert status payment: %s", err.Error())
				}
			} else {
				log.Fatalf("Error checking status payment: %s", err.Error())
			}
		} else {
			// update kalau ada perubahan
			if existing.Status != sp.Status {
				if err := s.db.Model(&existing).Update("status", sp.Status).Error; err != nil {
					log.Fatalf("Failed to update status payment: %s", err.Error())
				}
			}
		}
	}
	log.Println("Seeding status payment completed (sync)")
}
