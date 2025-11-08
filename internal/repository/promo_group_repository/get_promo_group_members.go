package promo_group_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (pgr *PromoGroupRepository) GetPromoGroupMembers(ctx context.Context, promoGroupID uint, limit, page int) ([]entity.User, int64, error) {
	db := pgr.db.GetTx(ctx)

	var members []model.User
	var total int64
	query := db.WithContext(ctx).
		Model(&model.User{}).
		Preload("AgentCompany").
		Select("users.id, users.full_name, users.agent_company_id").
		Where("promo_group_id = ?", promoGroupID)

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting promo groups", err.Error())
		return nil, total, err
	}

	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&members).Error; err != nil {
		logger.Error(ctx, "Error finding promo group", err.Error())
		return nil, total, err
	}

	membersEntity := make([]entity.User, 0, len(members))
	for _, m := range members {
		var e entity.User
		if err := utils.CopyPatch(&e, &m); err != nil {
			logger.Error(ctx, "Error copying user model to entity", err.Error())
			return nil, total, err
		}

		if m.AgentCompany != nil {
			e.AgentCompanyName = m.AgentCompany.Name
		}

		membersEntity = append(membersEntity, e)
	}
	return membersEntity, total, nil
}
