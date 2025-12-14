package booking_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) DeleteCartBooking(ctx context.Context, agentID uint, bookingDetailID uint) error {
	db := br.db.GetTx(ctx)

	// 🔍 Step 1: Ambil booking_id dari booking_detail dan validasi agent_id
	var booking model.Booking
	err := db.WithContext(ctx).
		Model(&model.Booking{}).
		Joins("JOIN booking_details bd ON bd.booking_id = bookings.id").
		Where("bd.id = ? AND bookings.agent_id = ?", bookingDetailID, agentID).
		Select("bookings.id").
		First(&booking).Error
	if err != nil {
		logger.Error(ctx, "booking detail not found or not owned by agent", err.Error())
		return fmt.Errorf("unauthorized or not found")
	}

	// 🗑️ Step 2: Hapus related records terlebih dahulu (cascade delete)
	// Delete BookingDetailAdditional records
	if err := db.WithContext(ctx).
		Unscoped().
		Where("booking_detail_id = ?", bookingDetailID).
		Delete(&model.BookingDetailAdditional{}).Error; err != nil {
		logger.Error(ctx, "failed to delete booking detail additionals", err.Error())
		return err
	}

	// Delete Invoice records if any
	if err := db.WithContext(ctx).
		Unscoped().
		Where("booking_detail_id = ?", bookingDetailID).
		Delete(&model.Invoice{}).Error; err != nil {
		logger.Error(ctx, "failed to delete invoice", err.Error())
		return err
	}

	// 🗑️ Step 3: Hapus booking_detail secara hard delete
	if err := db.WithContext(ctx).
		Unscoped().
		Where("id = ?", bookingDetailID).
		Delete(&model.BookingDetail{}).Error; err != nil {
		logger.Error(ctx, "failed to delete booking detail", err.Error())
		return err
	}

	// 🔍 Step 4: Cek apakah masih ada detail lain di booking yang sama
	var remaining int64
	if err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("booking_id = ?", booking.ID).
		Count(&remaining).Error; err != nil {
		logger.Error(ctx, "failed to count remaining booking details", err.Error())
		return err
	}

	// 🗑️ Step 5: Kalau kosong, hapus booking-nya juga
	if remaining == 0 {
		if err := db.WithContext(ctx).
			Unscoped().
			Where("id = ?", booking.ID).
			Delete(&model.Booking{}).Error; err != nil {
			logger.Error(ctx, "failed to delete empty booking cart", err.Error())
			return err
		}
	}

	return nil
}
