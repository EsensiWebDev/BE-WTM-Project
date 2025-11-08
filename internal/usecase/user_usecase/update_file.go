package user_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (uu *UserUsecase) UpdateFile(ctx context.Context, req *userdto.UpdateFileRequest) error {
	dataUser, err := uu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "Error getting user from context:", err.Error())
		return err
	}

	if dataUser == nil {
		logger.Error(ctx, "User not found in context")
		return errors.New("user not found in context")
	}

	userDB, err := uu.userRepo.GetUserByID(ctx, dataUser.ID)
	if err != nil {
		logger.Error(ctx, "Error getting user by Id:", err.Error())
		return err
	}

	if userDB == nil {
		logger.Error(ctx, "User not found in database")
		return errors.New("user not found in database")
	}

	var label, typeAccess string
	var assignTo *string
	switch req.FileType {
	case "photo":
		label = "selfie"
		assignTo = &userDB.PhotoSelfie
		typeAccess = constant.ConstPublic
	case "certificate":
		label = "certificate"
		assignTo = &userDB.Certificate
		typeAccess = constant.ConstPrivate
	case "name_card":
		label = "name_card"
		assignTo = &userDB.NameCard
		typeAccess = constant.ConstPrivate
	default:
		logger.Error(ctx, "Invalid file type:", req.FileType)
	}

	if err := uu.uploadAndAssign(ctx, userDB, req.File, label, assignTo, typeAccess); err != nil {
		logger.Error(ctx, "Error uploading and assigning user photo:", err.Error())
		return errors.New("error uploading and assigning user photo")
	}

	_, err = uu.userRepo.UpdateUser(ctx, userDB)
	if err != nil {
		logger.Error(ctx, "Error updating user photo:", err.Error())
		return err
	}

	return nil
}
