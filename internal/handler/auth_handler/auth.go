package auth_handler

import (
	"wtm-backend/config"
	"wtm-backend/internal/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
	config      *config.Config
}

func NewAuthHandler(auth domain.AuthUsecase, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authUsecase: auth,
		config:      config,
	}
}
