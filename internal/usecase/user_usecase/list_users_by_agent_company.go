package user_usecase

import (
	"context"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) ListUsersByAgentCompany(ctx context.Context, req *userdto.ListUsersByAgentCompanyRequest) (*userdto.ListUsersByAgentCompanyResponse, int64, error) {
	users, total, err := uu.userRepo.GetUsersByAgentCompany(ctx, req.ID, req.Search, req.Limit, req.Page)
	if err != nil {
		logger.Error(ctx, "Error while fetching users by agent company")
		return nil, total, err
	}

	resp := &userdto.ListUsersByAgentCompanyResponse{}
	respData := make([]userdto.ListUsersByAgentCompanyData, 0, len(users))
	for _, user := range users {
		respData = append(respData, userdto.ListUsersByAgentCompanyData{
			ID:   user.ID,
			Name: user.FullName,
		})
	}
	resp.Users = respData

	return resp, total, nil
}
