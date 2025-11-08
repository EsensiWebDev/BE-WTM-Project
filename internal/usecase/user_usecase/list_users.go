package user_usecase

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) ListUsers(ctx context.Context, req *userdto.ListUsersRequest) (*userdto.ListUsersResponse, error) {
	filterUser := filter.UserFilter{
		AgentCompanyID:    req.AgentCompanyID,
		PaginationRequest: req.PaginationRequest,
		Scope:             req.Scope,
	}
	if req.StatusID > 0 {
		filterUser.StatusID = &req.StatusID
	}

	if strings.TrimSpace(req.Role) != "" {
		roleID := getRoleID(req.Role)
		filterUser.RoleID = &roleID
	}

	users, total, err := uu.userRepo.GetUsers(ctx, filterUser)
	if err != nil {
		logger.Error(ctx, "Error getting users", err.Error())
		return nil, err
	}

	resp := &userdto.ListUsersResponse{
		Total: total,
	}
	resp.Users = make([]userdto.ListUserData, 0, len(users))
	for _, u := range users {

		var photoProfile string
		if strings.TrimSpace(u.PhotoSelfie) != "" {
			bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPublic)
			photoProfile, err = uu.fileStorage.GetFile(ctx, bucketName, u.PhotoSelfie)
			if err != nil {
				logger.Error(ctx, "Error getting user profile photo", err.Error())
				return nil, fmt.Errorf("failed to get user profile photo: %s", err.Error())
			}
		}

		var certificateURL string
		if strings.TrimSpace(u.Certificate) != "" {
			bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPrivate)
			certificateURL, err = uu.fileStorage.GetFile(ctx, bucketName, u.Certificate)
			if err != nil {
				logger.Error(ctx, "Error getting user certificate", err.Error())
				return nil, fmt.Errorf("failed to get user certificate: %s", err.Error())
			}
		}

		var nameCardURL string
		if strings.TrimSpace(u.NameCard) != "" {
			bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPrivate)
			nameCardURL, err = uu.fileStorage.GetFile(ctx, bucketName, u.NameCard)
			if err != nil {
				logger.Error(ctx, "Error getting user name card", err.Error())
				return nil, fmt.Errorf("failed to get user name card: %s", err.Error())
			}
		}

		var idCardURL string
		if strings.TrimSpace(u.PhotoIDCard) != "" {
			bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, constant.ConstPrivate)
			idCardURL, err = uu.fileStorage.GetFile(ctx, bucketName, u.PhotoIDCard)
			if err != nil {
				logger.Error(ctx, "Error getting user Id card", err.Error())
				return nil, fmt.Errorf("failed to get user Id card: %s", err.Error())
			}
		}

		data := userdto.ListUserData{
			ID:               u.ID,
			Name:             u.FullName,
			Email:            u.Email,
			PhoneNumber:      u.Phone,
			Status:           u.StatusName,
			PromoGroupName:   u.PromoGroupName,
			PromoGroupID:     u.PromoGroupID,
			AgentCompanyName: u.AgentCompanyName,
			KakaoTalkID:      u.KakaoTalkID,
			Photo:            photoProfile,
			Certificate:      certificateURL,
			NameCard:         nameCardURL,
			IdCard:           idCardURL,
		}
		if u.PromoGroupID != nil {
			data.PromoGroupID = u.PromoGroupID
		}
		resp.Users = append(resp.Users, data)
	}

	return resp, nil
}
