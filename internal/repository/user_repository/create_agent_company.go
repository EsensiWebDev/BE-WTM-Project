package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) CreateAgentCompany(ctx context.Context, agentCompany string) (*entity.AgentCompany, error) {
	db := ur.db.GetTx(ctx)

	modelAgentCompany := model.AgentCompany{
		Name: agentCompany,
	}

	err := db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", agentCompany).FirstOrCreate(&modelAgentCompany).Error
	if err != nil {
		logger.Error(ctx, "Error to add agent company", err.Error())
		return nil, err
	}

	var entityAgentCompany entity.AgentCompany
	if err := utils.CopyStrict(&entityAgentCompany, modelAgentCompany); err != nil {
		logger.Error(ctx, "Error copying agent company model to entity", err.Error())
		return nil, err
	}

	return &entityAgentCompany, nil
}
