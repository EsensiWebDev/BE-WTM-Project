package user_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) UpdateProfile(ctx context.Context, user *userdto.UpdateProfileRequest) error {
	userCtx, err := uu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error generating user from context", err.Error())
		return err
	}

	if userCtx == nil {
		logger.Warn(ctx,
			"User context is nil")
		return errors.New("user context is nil")
	}

	userDB, err := uu.userRepo.GetUserByID(ctx, userCtx.ID)
	if err != nil {
		logger.Error(ctx, "Error getting user by Id", err.Error())
		return err
	}

	if userDB == nil {
		logger.Warn(ctx,
			"User not found in database")
		return errors.New("user not found in database")
	}

	if user.FullName == userDB.FullName && user.Phone == userDB.Phone && user.Email == userDB.Email && user.KakaoTalkID == userDB.KakaoTalkID {
		logger.Info(ctx, "No changes detected in user profile", nil)
		return errors.New("no changes detected in user profile")
	}

	if userDB.RoleID == constant.RoleAgentID {
		if user.KakaoTalkID == "" {
			logger.Warn(ctx, "KakaoTalk ID is required for agents")
			return errors.New("kakaotalk ID is required for agents")
		}

		userDB.KakaoTalkID = user.KakaoTalkID
	}

	userDB.FullName = user.FullName
	userDB.Email = user.Email
	userDB.Phone = user.Phone

	_, err = uu.userRepo.UpdateUser(ctx, userDB)
	if err != nil {
		logger.Error(ctx, "Error updating user profile", err.Error())
		return err
	}

	return nil
}
