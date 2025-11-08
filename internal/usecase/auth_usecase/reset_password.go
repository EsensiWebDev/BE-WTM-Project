package auth_usecase

import (
	"context"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (au *AuthUsecase) ResetPassword(ctx context.Context, request *authdto.ResetPasswordRequest) error {
	return au.dbTrx.WithTransaction(ctx, func(nCtx context.Context) error {
		userID, err := au.authRepo.UsedTokenResetPassword(nCtx, request.Token)
		if err != nil {
			logger.Error(ctx, "Error marking reset password token as used:", err.Error())
			return err
		}

		encryptPass, err := utils.GeneratePassword(ctx, request.Password)
		if err != nil {
			logger.Error(ctx, "Error encrypting new password:", err.Error())
			return err
		}

		user, err := au.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			logger.Error(ctx, "Error fetching user by ID:", err.Error())
			return err
		}

		user.Password = encryptPass

		_, err = au.userRepo.UpdateUser(nCtx, user)
		if err != nil {
			logger.Error(ctx, "Error updating user password:", err.Error())
			return err
		}

		return nil
	})
}
