package auth_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/domain/entity"
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

	if user == nil {
		logger.Error(ctx, "User not found for email:", request.Email)
		return nil, fmt.Errorf("user not found for email: %s", request.Email)
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
		au.sendEmailNotification(newCtx, user.FullName, user.Email, token, user.RoleID, au.config.DurationLinkExpiration)
	}()

	return nil, nil
}

func (au *AuthUsecase) sendEmailNotification(ctx context.Context, name, email, token string, roleID uint, expiry time.Duration) {
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

	var url string
	if roleID == constant.RoleAgentID {
		url = au.config.URLFEAgent
	} else {
		url = au.config.URLFEAdmin
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", url, token)

	// Inject data
	data := EmailData{
		FullName:  name,
		ResetLink: resetLink,
		ExpiresIn: utils.HumanizeDuration(expiry),
	}

	subjectParsed := emailTemplate.Subject

	bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
	if err != nil {
		logger.Error(ctx, "Error parsing body HTML:", err.Error())
		return
	}

	bodyText := "Please view this email in HTML format." // Optional fallback

	emailTo := email

	emailLog := entity.EmailLog{
		To:              emailTo,
		Subject:         subjectParsed,
		Body:            bodyHTML,
		EmailTemplateID: uint(emailTemplate.ID),
	}
	metadataLog := entity.MetadataEmailLog{AgentName: name}
	emailLog.Meta = &metadataLog
	var dataEmail bool
	statusEmailID := constant.StatusEmailSuccessID
	if err = au.emailRepo.CreateEmailLog(ctx, &emailLog); err != nil {
		logger.Error(ctx, "Failed to create email log:", err)
		dataEmail = false
	} else {
		dataEmail = true
	}

	err = au.emailSender.Send(ctx, constant.ScopeAgent, emailTo, subjectParsed, bodyHTML, bodyText)
	if err != nil {
		logger.Error(ctx, "Error sending email:", err.Error())
		statusEmailID = constant.StatusEmailFailedID
		metadataLog.Notes = fmt.Sprintf("Error sending email: %s", err.Error())
		emailLog.Meta = &metadataLog
	}

	if dataEmail {
		emailLog.StatusID = uint(statusEmailID)
		if err := au.emailRepo.UpdateStatusEmailLog(ctx, &emailLog); err != nil {
			logger.Error(ctx, "Failed to update email log:", err.Error())
		}
	}

}

type EmailData struct {
	FullName  string
	ResetLink string
	ExpiresIn string // e.g. "30 minutes"
}
