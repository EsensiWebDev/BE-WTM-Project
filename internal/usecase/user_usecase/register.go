package user_usecase

import (
	"context"
	"errors"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (uu *UserUsecase) Register(ctx context.Context, userReq *userdto.RegisterRequest) error {
	passCrypt, err := utils.GeneratePassword(ctx, userReq.Password)
	if err != nil {
		logger.Error(ctx, "Error generating password", err.Error())
		return err
	}

	user := &entity.User{
		FullName:    userReq.FullName,
		Username:    userReq.Username,
		Password:    passCrypt,
		StatusID:    constant.DefaultStatusSign,
		RoleID:      constant.DefaultRoleAgent,
		Email:       userReq.Email,
		Phone:       userReq.Phone,
		KakaoTalkID: userReq.KakaoTalkID,
	}

	return uu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {
		if strings.TrimSpace(userReq.AgentCompany) != "" {

			agentCompany, err := uu.userRepo.CreateAgentCompany(txCtx, userReq.AgentCompany)
			if err != nil {
				logger.Error(ctx, "Error to add agent company", err.Error())
				return err
			}

			user.AgentCompanyID = &agentCompany.ID
		}

		userDB, err := uu.userRepo.CreateUser(txCtx, user)
		if err != nil {
			logger.Error(ctx, "Error to add user", err.Error())
			return err
		}

		if err := uu.uploadAndAssign(txCtx, userDB, userReq.PhotoSelfie, "selfie", &userDB.PhotoSelfie, constant.ConstPublic); err != nil {
			logger.Error(ctx, "Error uploading selfie photo", err.Error())
			return errors.New("upload selfie photo error")
		}

		if err := uu.uploadAndAssign(txCtx, userDB, userReq.PhotoIDCard, "id_card", &userDB.PhotoIDCard, constant.ConstPrivate); err != nil {
			logger.Error(ctx, "Error uploading Id card photo", err.Error())
			return err
		}

		if userReq.Certificate != nil {
			if err := uu.uploadAndAssign(txCtx, userDB, userReq.Certificate, "certificate", &userDB.Certificate, constant.ConstPrivate); err != nil {
				logger.Error(ctx, "Error uploading certificate", err.Error())
				return err
			}
		}

		if err := uu.uploadAndAssign(txCtx, userDB, userReq.NameCard, "name_card", &userDB.NameCard, constant.ConstPrivate); err != nil {
			logger.Error(ctx, "Error uploading name card photo", err.Error())
			return err
		}

		_, err = uu.userRepo.UpdateUser(txCtx, userDB)
		if err != nil {
			logger.Error(ctx, "Error updating user after upload", err.Error())
			return err
		}

		return nil
	})
}
