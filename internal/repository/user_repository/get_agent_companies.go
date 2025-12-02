package user_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetAgentCompanies(ctx context.Context, search string, limit, page int) ([]entity.AgentCompany, int64, error) {
	db := ur.db.GetTx(ctx)

	var modelAgentCompany []model.AgentCompany
	var total int64

	query := db.WithContext(ctx).
		Select("id, name").
		Model(&model.AgentCompany{})

	if strings.TrimSpace(search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(search)
		query = query.Where("name ILIKE ? ", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting agent companies", err.Error())
		return nil, total, err
	}

	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&modelAgentCompany).Error; err != nil {
		logger.Error(ctx, "Error to get agent companies", err.Error())
		return nil, total, err
	}

	var entityAgentCompany []entity.AgentCompany
	if err := utils.CopyPatch(&entityAgentCompany, modelAgentCompany); err != nil {
		logger.Error(ctx, "Error copying agent company model to entity", err.Error())
		return nil, total, err
	}

	return entityAgentCompany, total, nil

}
