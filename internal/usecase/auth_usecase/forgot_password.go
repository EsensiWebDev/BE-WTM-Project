package auth_usecase

import (
	"context"
	"time"
	dtoauth "wtm-backend/internal/dto/authdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (au *AuthUsecase) ForgotPassword(ctx context.Context, request *dtoauth.ForgotPasswordRequest) (*dtoauth.ForgotPasswordResponse, error) {
	// Validate if a recent password reset request has been made
	exists, expiresAt, err := au.authRepo.ValidateEmailForgotPassword(ctx, request.Email)
	if err != nil {
		logger.Error(ctx, "Error validating email forgot password:", err.Error())
		return nil, err
	}
	if exists {
		logger.Warn(ctx,
			"Password reset request already exists for email:", request.Email)
		expTime := time.Now().Add(expiresAt).Format(time.RFC3339)
		resp := &dtoauth.ForgotPasswordResponse{
			ExpiresAt: expTime,
		}
		return resp, nil

	}

	durationExpiration := au.config.DurationLinkExpiration
	// Validate if the email exists in the system
	user, err := au.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		logger.Error(ctx, "Error getting user by email:", err.Error())
		return nil, err
	}

	token, err := utils.GenerateSecureToken()
	if err != nil {
		logger.Error(ctx, "Error generating secure token:", err.Error())
		return nil, err
	}
	hashed := utils.HashToken(token)
	// Store the hashed token in the database
	if err := au.authRepo.CreatePasswordResetToken(ctx, user.ID, hashed, durationExpiration); err != nil {
		logger.Error(ctx, "Error creating password reset token:", err.Error())
		return nil, err
	}

	// Set a flag in Redis to limit the frequency of password reset requests
	if err := au.authRepo.SetEmailForgotPassword(ctx, request.Email, durationExpiration); err != nil {
		logger.Error(ctx, "Error setting email forgot password in redis:", err.Error())
		return nil, err
	}

	// Send the password reset email
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), au.config.DurationCtxTOSlow)
		defer cancel()
		au.sendEmailNotification(newCtx, user.FullName, user.Email, token, au.config.DurationLinkExpiration)
	}()

	return nil, nil
}

func (au *AuthUsecase) sendEmailNotification(ctx context.Context, name, email, token string, expiry time.Duration) {
	var statusEmail = constant.EmailForgotPassword

	emailTemplate, err := au.emailRepo.GetEmailTemplateByName(ctx, statusEmail)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for status:", statusEmail)
		return
	}

	resetLink := "https://yourdomain.com/reset-password?token=" + token

	// Inject data
	data := EmailData{
		FullName:  name,
		ResetLink: resetLink,
		ExpiresIn: utils.HumanizeDuration(expiry),
	}

	bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
	if err != nil {
		logger.Error(ctx, "Error parsing body HTML:", err.Error())
		return
	}

	bodyText := "Please view this email in HTML format." // Optional fallback

	err = au.emailSender.Send(ctx, email, emailTemplate.Subject, bodyHTML, bodyText)
	if err != nil {
		logger.Error(ctx, "Error sending email:", err.Error())
	}

}

type EmailData struct {
	FullName  string
	ResetLink string
	ExpiresIn string // e.g. "30 minutes"
}
