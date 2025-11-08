package domain

import (
	"context"
	"time"
	"wtm-backend/internal/dto/authdto"
)

type AuthUsecase interface {
	Login(ctx context.Context, request *authdto.LoginRequest) (*authdto.LoginResponse, string, error)
	RefreshToken(ctx context.Context, token string) (*authdto.LoginResponse, error)
	Logout(ctx context.Context) error
	ForgotPassword(ctx context.Context, request *authdto.ForgotPasswordRequest) (*authdto.ForgotPasswordResponse, error)
	ValidateTokenResetPassword(ctx context.Context, req *authdto.ValidateTokenResetPasswordRequest) (*authdto.ValidateTokenResetPasswordResponse, error)
	ResetPassword(ctx context.Context, request *authdto.ResetPasswordRequest) error
}

type AuthRepository interface {
	SetAccessToken(ctx context.Context, userID uint, accessToken string, expiry time.Duration) error
	ValidateAccessToken(ctx context.Context, userID uint, accessToken string) (bool, error)
	DeleteAccessToken(ctx context.Context, userID uint) error
	CreatePasswordResetToken(ctx context.Context, userID uint, token string, expiry time.Duration) error
	FindActiveResetTokenByUserID(ctx context.Context, userID uint) (string, error)
	FindActiveResetTokenByToken(ctx context.Context, token string) (string, error)
	SetEmailForgotPassword(ctx context.Context, email string, duration time.Duration) error
	ValidateEmailForgotPassword(ctx context.Context, email string) (bool, time.Duration, error)
	DeleteEmailForgotPassword(ctx context.Context, email string) error
	UsedTokenResetPassword(ctx context.Context, token string) (uint, error)
}
