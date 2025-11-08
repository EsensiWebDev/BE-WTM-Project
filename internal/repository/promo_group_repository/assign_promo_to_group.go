package promo_group_repository

import (
	"context"
	"fmt"
	"wtm-backend/internal/infrastructure/database/model"
)

func (pgr *PromoGroupRepository) AssignPromoToGroup(ctx context.Context, promoGroupID uint, promoID uint) error {
	db := pgr.db.GetTx(ctx)

	// 1️⃣ Ambil promo group
	var group model.PromoGroup
	if err := db.WithContext(ctx).First(&group, promoGroupID).Error; err != nil {
		return fmt.Errorf("promo group not found: %w", err)
	}

	// 2️⃣ Ambil promo
	var promo model.Promo
	if err := db.WithContext(ctx).First(&promo, promoID).Error; err != nil {
		return fmt.Errorf("promo not found: %w", err)
	}

	// 3️⃣ Assign promo ke group (many2many append)
	if err := db.WithContext(ctx).Model(&group).Association("Promos").Append(&promo); err != nil {
		return fmt.Errorf("failed to assign promo to group: %w", err)
	}

	return nil

}
