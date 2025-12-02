package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/jwt"
	"wtm-backend/pkg/logger"
)

type userContextKey struct{}

const (
	roleKey       = "role"
	permissionKey = "permissions"
)

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn(ctx, "Authorization header is not bearer token")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseToken(ctx, tokenStr, m.jwtSecret)
		if err != nil {
			logger.Error(ctx, "Error to parse token", err.Error())
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if claims.Type != "access" {
			logger.Warn(ctx, "Invalid type token")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user := m.GenerateUserFromClaimToken(claims)

		// Cek apakah user ada di redis
		ok, err := m.authRepo.ValidateAccessToken(c.Request.Context(), user.ID, tokenStr)
		if err != nil {
			logger.Error(ctx, "Failed to validate access token", "valid", ok, "err", err.Error())
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		} else if !ok {
			logger.Error(ctx, "Access token not valid", "valid", ok)
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Simpan ke context.Context
		ctx = context.WithValue(ctx, userContextKey{}, user)

		// Replace context di request agar bisa dipakai downstream
		c.Request = c.Request.WithContext(ctx)

		c.Set(roleKey, claims.User.Role)
		c.Set(permissionKey, claims.User.Permissions)

		c.Next()
	}
}

func (m *Middleware) GenerateUserFromClaimToken(claims *jwt.JwtClaims) *entity.User {
	return &entity.User{
		ID:          claims.User.ID,
		Username:    claims.User.Username,
		RoleName:    claims.User.Role,
		RoleID:      claims.User.RoleID,
		Permissions: claims.User.Permissions,
		PhotoSelfie: claims.User.PhotoURL,
		FullName:    claims.User.FullName,
	}
}

func (m *Middleware) RequirePermission(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		role, exists := c.Get(roleKey)
		if !exists {
			logger.Warn(ctx, "Role not found in context")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if role == constant.RoleSuperAdminCap {
			// Superadmin has all permissions, skip check
			c.Next()
			return
		}

		rawPerms, exists := c.Get(permissionKey)
		if !exists {
			logger.Warn(ctx, "Permissions not found in context")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		perms, ok := rawPerms.([]string)
		if !ok {
			logger.Warn(ctx, "Permissions not found in context")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		for _, p := range perms {
			if p == required {
				c.Next()
				return
			}
		}

		logger.Warn(ctx, "Required permission not found")
		response.Error(c, http.StatusForbidden, "Forbidden")
		c.Abort()
		return
	}
}

func (m *Middleware) RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		role, exists := c.Get(roleKey)
		if !exists {
			logger.Warn(ctx, "Role not found in context")
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if role == constant.RoleSuperAdmin {
			c.Next()
			return
		}

		if role == required {
			c.Next()
			return
		}

		logger.Warn(ctx, "Required role not found")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
	}
}

func (m *Middleware) GenerateUserFromContext(c context.Context) (*entity.User, error) {
	user, exists := c.Value(userContextKey{}).(*entity.User)
	if !exists || user == nil {
		logger.Warn(c,
			"User not found in context")
		return nil, fmt.Errorf("user not found in context")
	}

	return user, nil
}
