package user_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUsersByAgentCompany(ctx context.Context, agentCompanyID uint, search string, limit, page int) ([]entity.User, int64, error) {
	db := ur.db.GetTx(ctx)

	var user []model.User
	var total int64
	query := db.WithContext(ctx).
		Select("id, full_name").
		Model(&model.User{}).
		Where("agent_company_id = ?", agentCompanyID)

	if strings.TrimSpace(search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(search)
		query = query.Where("LOWER(full_name) ILIKE ?", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting users by agent company", err.Error())
		return nil, total, err
	}

	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&user).Error; err != nil {
		logger.Error(ctx, "Error to get user by agent company", err.Error())
		return nil, total, err
	}

	var entityUser []entity.User
	if err := utils.CopyPatch(&entityUser, user); err != nil {
		logger.Error(ctx, "Error copying user model to entity", err.Error())
		return nil, total, err
	}

	return entityUser, total, nil
}
