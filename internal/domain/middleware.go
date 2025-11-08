package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/jwt"
)

type Middleware interface {
	GenerateUserFromContext(ctx context.Context) (*entity.User, error)
	GenerateUserFromClaimToken(claim *jwt.JwtClaims) *entity.User
}
