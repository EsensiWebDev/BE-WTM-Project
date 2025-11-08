package user_usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

type UserUsecase struct {
	userRepo       domain.UserRepository
	authRepo       domain.AuthRepository
	promoGroupRepo domain.PromoGroupRepository
	emailRepo      domain.EmailRepository
	config         *config.Config
	fileStorage    domain.StorageClient
	middleware     domain.Middleware
	dbTrx          domain.DatabaseTransaction
	emailSender    domain.EmailSender
}

func NewUserUsecase(userRepo domain.UserRepository, authRepo domain.AuthRepository, promoGroupRepo domain.PromoGroupRepository, emailRepo domain.EmailRepository, config *config.Config, minio domain.StorageClient, middleware domain.Middleware, dbTrx domain.DatabaseTransaction, emailSender domain.EmailSender) *UserUsecase {
	return &UserUsecase{
		userRepo:       userRepo,
		authRepo:       authRepo,
		promoGroupRepo: promoGroupRepo,
		emailRepo:      emailRepo,
		config:         config,
		fileStorage:    minio,
		middleware:     middleware,
		dbTrx:          dbTrx,
		emailSender:    emailSender,
	}
}

func (uu *UserUsecase) uploadFile(ctx context.Context, userID uint, file *multipart.FileHeader, prefix string, typeAccess string) (string, error) {
	if file == nil {
		return "", nil
	}
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			logger.Error(ctx, "failed to close file", err.Error())
		}
	}(f)

	bucketName := fmt.Sprintf("%s-%s", constant.ConstUser, typeAccess)
	filename := fmt.Sprintf("%s_%d_%d%s", prefix, userID, time.Now().Unix(), filepath.Ext(file.Filename))
	return uu.fileStorage.UploadFile(ctx, f, file, bucketName, filename)
}

func (uu *UserUsecase) uploadAndAssign(ctx context.Context, user *entity.User, file *multipart.FileHeader, label string, assign *string, typeAccess string) error {
	url, err := uu.uploadFile(ctx, user.ID, file, label, typeAccess)
	if err != nil {
		logger.Error(ctx, "Error to upload file", err.Error())
		return err
	}
	*assign = url
	return nil
}

func getRoleID(role string) uint {

	role = strings.ToLower(role)
	role = strings.TrimSpace(role)
	role = strings.ReplaceAll(role, " ", "_")

	switch role {
	case constant.RoleSuperAdmin:
		return constant.RoleSuperAdminID
	case constant.RoleAdmin:
		return constant.RoleAdminID
	case constant.RoleAgent:
		return constant.RoleAgentID
	case constant.RoleSupport:
		return constant.RoleSupportID
	default:
		return 0
	}

}

func getStatusID(isActive bool) uint {
	if isActive {
		return constant.StatusUserActiveID // Active
	}
	return constant.StatusUserWaitingApprovalID // Sign
}
