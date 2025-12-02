package booking_repository

import (
	"context"
	"gorm.io/gorm"
	"wtm-backend/internal/infrastructure/database/model"
)

func (br *BookingRepository) RemoveGuestsFromCart(ctx context.Context, agentID uint, bookingID uint, guests []string) error {
	db := br.db.GetTx(ctx)

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Step 1: Validasi booking ID dan agent ID
		var exists bool
		if err := tx.Model(&model.Booking{}).
			Select("count(*) > 0").
			Where("id = ? AND agent_id = ?", bookingID, agentID).
			Find(&exists).Error; err != nil {
			return err
		}
		if !exists {
			return nil // Booking tidak ditemukan, tidak ada yang dihapus
		}

		// Step 2: Hapus guest dari tabel booking_guests
		if err := tx.Unscoped().
			Where("booking_id = ? AND name IN (?)", bookingID, guests).
			Delete(&model.BookingGuest{}).Error; err != nil {
			return err
		}

		// Step 3: Update booking_detail yang punya guest tersebut
		// Asumsi field "guest" di BookingDetail menyimpan nama guest (string)
		if err := tx.Model(&model.BookingDetail{}).
			Where("booking_id = ? AND guest IN (?)", bookingID, guests).
			Update("guest", "").Error; err != nil {
			return err
		}

		return nil
	})
}
