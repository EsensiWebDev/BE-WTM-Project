package user_usecase

import (
	"context"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) ListStatusUsers(ctx context.Context, req *userdto.ListStatusUsersRequest) (*userdto.ListStatusUsersResponse, int64, error) {
	filterRepo := &filter.DefaultFilter{}
	filterRepo.PaginationRequest = req.PaginationRequest

	statusUsers, total, err := uu.userRepo.GetStatusUsers(ctx, filterRepo)
	if err != nil {
		logger.Error(ctx, "Error getting status users:", err.Error())
		return nil, 0, err
	}

	resp := &userdto.ListStatusUsersResponse{
		StatusUsers: statusUsers,
	}

	return resp, total, nil

}
