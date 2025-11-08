package middleware

import (
	"time"
	"wtm-backend/config"
	"wtm-backend/internal/domain"
	"wtm-backend/internal/repository/auth_repository"
)

type Middleware struct {
	authRepo       domain.AuthRepository
	jwtSecret      string
	maxAgeCors     time.Duration
	allowedOrigins []string
}

func NewMiddleware(config *config.Config, authRepo *auth_repository.AuthRepository) *Middleware {
	return &Middleware{
		jwtSecret: config.JWTSecret,
		authRepo:  authRepo,
	}
}
