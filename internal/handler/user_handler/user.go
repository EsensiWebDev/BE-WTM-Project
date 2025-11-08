package user_handler

import (
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/usecase/user_usecase"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
	config      *config.Config
}

func NewUserHandler(user *user_usecase.UserUsecase, config *config.Config) *UserHandler {
	return &UserHandler{
		userUsecase: user,
		config:      config,
	}
}
