package user_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/dto/userdto"
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

	userDB.FullName = req.FullName
	userDB.Email = req.Email
	userDB.Phone = req.Phone
	userDB.StatusID = getStatusID(req.IsActive)

	_, err = uu.userRepo.UpdateUser(ctx, userDB)
	if err != nil {
		logger.Error(ctx, "Error updating user profile", err.Error())
		return err
	}

	return nil
}
