package user_usecase

import (
	"context"
	"errors"
	"strings"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) UpdateUserByAdmin(ctx context.Context, req *userdto.UpdateUserByAdminRequest) error {

	userDB, err := uu.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		logger.Error(ctx, "Error getting user by Id", err.Error())
		return err
	}

	if userDB == nil {
		logger.Warn(ctx,
			"User not found in database")
		return errors.New("user not found in database")
	}

	var isNeedLogout bool

	if userDB.Username != req.Username {
		user, err := uu.userRepo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			logger.Error(ctx, "Error getting user by username", err.Error())
			return err
		}
		if user != nil && user.ID != userDB.ID {
			logger.Warn(ctx, "User with username already exists")
			return errors.New("user with username already exists")
		}
		isNeedLogout = true
	}

	if userDB.Email != req.Email {
		user, err := uu.userRepo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			logger.Error(ctx, "Error getting user by email", err.Error())
			return err
		}
		if user != nil && user.ID != userDB.ID {
			logger.Warn(ctx, "User with email already exists")
			return errors.New("user with email already exists")
		}
		isNeedLogout = true
	}

	if userDB.Phone != req.Phone {
		user, err := uu.userRepo.GetUserByPhone(ctx, req.Phone)
		if err != nil {
			logger.Error(ctx, "Error getting user by phone", err.Error())
			return err
		}
		if user != nil && user.ID != userDB.ID {
			logger.Warn(ctx, "User with phone already exists")
			return errors.New("user with phone already exists")
		}
	}

	trxErr := uu.dbTrx.WithTransaction(ctx, func(txCtx context.Context) error {

		if userDB.RoleID == constant.RoleAgentID {
			if req.PromoGroupID > 0 {
				promoGroup, err := uu.promoGroupRepo.GetPromoGroupByID(txCtx, req.PromoGroupID)
				if err != nil {
					logger.Error(txCtx, "Error to get promo group by Id", err.Error())
					return err
				}

				if promoGroup != nil && promoGroup.ID > 0 {
					userDB.PromoGroupID = &req.PromoGroupID
				}
			}

			if strings.TrimSpace(req.AgentCompany) != "" {

				agentCompany, err := uu.userRepo.CreateAgentCompany(txCtx, req.AgentCompany)
				if err != nil {
					logger.Error(txCtx, "Error to add agent company", err.Error())
					return err
				}

				userDB.AgentCompanyID = &agentCompany.ID
			} else {
				userDB.AgentCompanyID = nil
			}

			if req.PhotoSelfie != nil {
				if err := uu.uploadAndAssign(txCtx, userDB, req.PhotoSelfie, "selfie", &userDB.PhotoSelfie, constant.ConstPublic); err != nil {
					logger.Error(txCtx, "Error uploading selfie photo", err.Error())
					return errors.New("upload selfie photo error")
				}
			}

			if req.PhotoIDCard != nil {
				if err := uu.uploadAndAssign(txCtx, userDB, req.PhotoIDCard, "id_card", &userDB.PhotoIDCard, constant.ConstPrivate); err != nil {
					logger.Error(txCtx, "Error uploading Id card photo", err.Error())
					return err
				}
			}

			if req.Certificate != nil {
				if err := uu.uploadAndAssign(txCtx, userDB, req.Certificate, "certificate", &userDB.Certificate, constant.ConstPrivate); err != nil {
					logger.Error(txCtx, "Error uploading certificate", err.Error())
					return err
				}
			}

			if req.NameCard != nil {
				if err := uu.uploadAndAssign(txCtx, userDB, req.NameCard, "name_card", &userDB.NameCard, constant.ConstPrivate); err != nil {
					logger.Error(txCtx, "Error uploading name card photo", err.Error())
					return err
				}
			}
		}

		userDB.FullName = req.FullName
		userDB.Email = req.Email
		userDB.Phone = req.Phone
		userDB.Username = req.Username
		userDB.KakaoTalkID = req.KakaoTalkID
		userDB.StatusID = getStatusID(req.IsActive)
		if userDB.StatusID == constant.StatusUserInactiveID {
			isNeedLogout = true
		}

		_, err = uu.userRepo.UpdateUser(txCtx, userDB)
		if err != nil {
			logger.Error(txCtx, "Error updating user profile", err.Error())
			return err
		}

		if isNeedLogout {
			if err := uu.authRepo.DeleteAccessToken(txCtx, userDB.ID); err != nil {
				logger.Error(txCtx, "Error to delete access token", err.Error())
				return errors.New("failed to delete access token")
			}
		}
		return nil
	})

	if trxErr != nil {
		logger.Error(ctx, "Error to update user profile", trxErr.Error())
		return trxErr
	}

	return nil
}
