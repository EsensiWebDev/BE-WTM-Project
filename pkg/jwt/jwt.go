package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/pkg/logger"
)

type JwtClaims struct {
	User *entity.UserMin `json:"user"`
	Type string          `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// Untuk akses token
func GenerateAccessToken(user *entity.User, secret string, expiry time.Duration) (string, error) {
	return generateJWT(user, "access", secret, expiry)
}

// Untuk refresh token
func GenerateRefreshToken(user *entity.User, secret string, expiry time.Duration) (string, error) {
	return generateJWT(user, "refresh", secret, expiry)
}

// ParseToken memverifikasi token dan mengembalikan klaim jika valid
func ParseToken(ctx context.Context, tokenStr, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warn(ctx, "Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		logger.Error(ctx, "Error parsing token", err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		logger.Warn(ctx,
			"Invalid token claims or token is not valid")
		return nil, fmt.Errorf("invalid token")
	}

	// Validasi expired secara eksplisit
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		logger.Warn(ctx,
			"Token expired at %v", claims.ExpiresAt.Time)
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

// generateJWT adalah fungsi umum untuk membuat token JWT
func generateJWT(user *entity.User, tokenType string, secret string, expiry time.Duration) (string, error) {
	claims := JwtClaims{
		User: &entity.UserMin{
			ID:          user.ID,
			Username:    user.Username,
			Role:        user.RoleName,
			RoleID:      user.RoleID,
			Permissions: user.Permissions,
			PhotoURL:    user.PhotoSelfie,
			FullName:    user.FullName,
		},
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
