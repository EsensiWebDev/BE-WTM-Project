package auth_usecase

import (
	"context"
	"wtm-backend/internal/dto/authdto"
	"wtm-backend/pkg/logger"
)

func (au *AuthUsecase) ValidateTokenResetPassword(ctx context.Context, req *authdto.ValidateTokenResetPasswordRequest) (*authdto.ValidateTokenResetPasswordResponse, error) {

	resp := &authdto.ValidateTokenResetPasswordResponse{}
	email, err := au.authRepo.FindActiveResetTokenByToken(ctx, req.Token)
	if err != nil {
		logger.Error(ctx, "Error validating forgot password token:", err.Error())
		return nil, err
	}
	resp.Email = email

	return resp, nil
}
