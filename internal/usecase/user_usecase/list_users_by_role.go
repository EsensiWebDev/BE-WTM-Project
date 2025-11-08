package user_usecase

import (
	"context"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) ListUsersByRole(ctx context.Context, req *userdto.ListUsersByRoleRequest) (*userdto.ListUsersByRoleResponse, int64, error) {
	roleID := getRoleID(req.Role)
	dataUser, total, err := uu.userRepo.GetUserByRole(ctx, roleID, req.Search, req.Limit, req.Page)
	if err != nil {
		logger.Error(ctx, "GetUserByRole failed", err.Error())
		return nil, total, err
	}

	resp := &userdto.ListUsersByRoleResponse{}

	respData := make([]userdto.ListUsersByRoleData, 0, len(dataUser))
	for _, user := range dataUser {
		data := userdto.ListUsersByRoleData{
			ID:               user.ID,
			Name:             user.FullName,
			Email:            user.Email,
			PhoneNumber:      user.Phone,
			Status:           user.StatusName,
			PromoGroupName:   user.PromoGroupName,
			AgentCompanyName: user.AgentCompanyName,
			KakaoTalkID:      user.KakaoTalkID,
		}
		if user.PromoGroupID != nil {
			data.PromoGroupID = *user.PromoGroupID
		}
		respData = append(respData, data)
	}

	resp.User = respData

	return resp, total, nil
}
