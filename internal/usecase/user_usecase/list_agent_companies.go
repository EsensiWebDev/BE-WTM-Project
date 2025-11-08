package user_usecase

import (
	"context"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) ListAgentCompanies(ctx context.Context, req *userdto.ListAgentCompaniesRequest) (*userdto.ListAgentCompaniesResponse, int64, error) {
	agentCompanies, total, err := uu.userRepo.GetAgentCompanies(ctx, req.Search, req.Limit, req.Page)
	if err != nil {
		logger.Error(ctx, "Error to get agent companies", err.Error())
		return nil, total, err
	}

	resp := &userdto.ListAgentCompaniesResponse{
		AgentCompanies: agentCompanies,
	}

	return resp, total, nil
}
