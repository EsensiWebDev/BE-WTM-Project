package booking_repository

import (
	"context"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/infrastructure/database/model"

	"gorm.io/gorm"
)

func (br *BookingRepository) RemoveGuestsFromCart(ctx context.Context, agentID uint, bookingID uint, guests []bookingdto.GuestInfo) error {
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

		// Step 2: Collect guest names for updating booking_detail.guest field
		var guestNames []string

		// Step 3: Delete guests using composite key (booking_id, name, honorific, category, age)
		for _, guest := range guests {
			query := tx.Unscoped().
				Where("booking_id = ?", bookingID).
				Where("name = ?", guest.Name).
				Where("honorific = ?", guest.Honorific).
				Where("category = ?", guest.Category)

			// Handle age: if category is Child, age must match; if Adult, age should be NULL
			if guest.Category == "Child" && guest.Age != nil {
				query = query.Where("age = ?", *guest.Age)
			} else {
				query = query.Where("age IS NULL")
			}

			// Delete the guest
			if err := query.Delete(&model.BookingGuest{}).Error; err != nil {
				return err
			}

			// Collect name for updating booking_detail
			guestNames = append(guestNames, guest.Name)
		}

		// Step 4: Update booking_detail yang punya guest tersebut
		// Field "guest" di BookingDetail menyimpan nama guest (string), jadi kita perlu update berdasarkan nama
		if len(guestNames) > 0 {
			if err := tx.Model(&model.BookingDetail{}).
				Where("booking_id = ? AND guest IN (?)", bookingID, guestNames).
				Update("guest", "").Error; err != nil {
				return err
			}
		}

		return nil
	})
}
