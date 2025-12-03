package booking_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) UpdateDetailBookingDetail(ctx context.Context, bookingDetailID uint, room *entity.DetailRoom, promo *entity.DetailPromo, price float64, additionals []entity.BookingDetailAdditional) error {
	db := br.db.GetTx(ctx)

	updates := make(map[string]interface{})

	if price != 0 {
		updates["price"] = price
	}

	if room != nil {
		detailRoom, err := json.Marshal(room)
		if err != nil {
			return fmt.Errorf("failed to marshal room details: %w", err)
		}
		updates["detail_room"] = detailRoom
	}

	if promo != nil {
		detailPromo, err := json.Marshal(promo)
		if err != nil {
			return fmt.Errorf("failed to marshal promo details: %w", err)
		}
		updates["detail_promo"] = detailPromo
	}

	// Jika tidak ada yang diupdate, langsung return
	if len(updates) == 0 {
		return nil
	}

	updates["updated_at"] = gorm.Expr("NOW()")

	// Eksekusi single update query
	err := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Where("id = ?", bookingDetailID).
		Updates(updates).
		Error

	if err != nil {
		logger.Error(ctx, "failed to update booking detail: ", err.Error())
		return fmt.Errorf("failed to update booking detail: %w", err)
	}

	// Jika ada yang diupdate, update detail_additional
	if len(additionals) > 0 {
		for _, additional := range additionals {
			if err := db.WithContext(ctx).
				Model(&model.BookingDetailAdditional{}).
				Where("booking_detail_id = ?", bookingDetailID).
				Where("id = ?", additional.ID).
				Updates(map[string]interface{}{
					"price":           additional.Price,
					"name_additional": additional.NameAdditional,
				}).Error; err != nil {
				logger.Error(ctx, "failed to update booking detail additional: ", err.Error())
				return fmt.Errorf("failed to update booking detail additional: %w", err)
			}
		}
	}

	return nil
}
