package banner_repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (br *BannerRepository) UpdateOrderBanner(ctx context.Context, id string, direction string) error {
	db := br.db.GetTx(ctx)

	// Ambil banner sekarang
	var current model.Banner
	if err := db.WithContext(ctx).
		Where("external_id = ?", id).
		First(&current).Error; err != nil {
		logger.Error(ctx, "Error getting current banner", err.Error())
		return err
	}

	// Cari banner tetangga
	var neighbor model.Banner
	switch strings.TrimSpace(strings.ToLower(direction)) {
	case "up":
		// Cari banner dengan order lebih besar (di atas), paling kecil di atasnya
		if err := db.WithContext(ctx).
			Where("display_order > ?", current.DisplayOrder).
			Order("display_order ASC").
			First(&neighbor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cannot move up, already at top")
			}
			return err
		}

	case "down":
		// Cari banner dengan order lebih kecil (di bawah), paling besar di bawahnya
		if err := db.WithContext(ctx).
			Where("display_order < ?", current.DisplayOrder).
			Order("display_order DESC").
			First(&neighbor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cannot move down, already at bottom")
			}
			return err
		}

	default:
		return fmt.Errorf("invalid direction: %s", direction)
	}

	// Swap posisi
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update neighbor ke posisi current
		if err := tx.Model(&model.Banner{}).
			Where("id = ?", neighbor.ID).
			Update("display_order", current.DisplayOrder).Error; err != nil {
			return err
		}

		// Update current ke posisi neighbor
		if err := tx.Model(&model.Banner{}).
			Where("id = ?", current.ID).
			Update("display_order", neighbor.DisplayOrder).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(ctx, "Error swapping banner order", err.Error())
		return err
	}

	return nil
}
