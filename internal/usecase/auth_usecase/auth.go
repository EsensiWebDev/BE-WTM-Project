package auth_usecase

import (
	"wtm-backend/config"
	"wtm-backend/internal/domain"
)

type AuthUsecase struct {
	userRepo    domain.UserRepository
	authRepo    domain.AuthRepository
	middleware  domain.Middleware
	config      *config.Config
	fileStorage domain.StorageClient
	emailSender domain.EmailSender
	emailRepo   domain.EmailRepository
	dbTrx       domain.DatabaseTransaction
}

func NewAuthUsecase(
	userRepo domain.UserRepository,
	authRepo domain.AuthRepository,
	config *config.Config,
	fileStorage domain.StorageClient,
	middleware domain.Middleware,
	emailSender domain.EmailSender,
	emailRepo domain.EmailRepository,
	dbTrx domain.DatabaseTransaction,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		authRepo:    authRepo,
		middleware:  middleware,
		config:      config,
		fileStorage: fileStorage,
		emailSender: emailSender,
		emailRepo:   emailRepo,
		dbTrx:       dbTrx,
	}
}
