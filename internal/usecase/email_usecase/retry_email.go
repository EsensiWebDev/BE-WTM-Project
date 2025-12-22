package email_usecase

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) RetryEmail(ctx context.Context, req *emaildto.RetryEmailRequest) (*emaildto.RetryEmailResponse, error) {
	// Get email log
	emailLog, err := eu.emailRepo.GetEmailLogByID(ctx, req.ID)
	if err != nil {
		logger.Error(ctx, "failed to get email log for retry", err.Error())
		return &emaildto.RetryEmailResponse{
			Success: false,
			Message: "Email log not found",
		}, err
	}

	// Check if email is failed
	if emailLog.StatusID != constant.StatusEmailFailedID {
		return &emaildto.RetryEmailResponse{
			Success: false,
			Message: "Only failed emails can be retried",
		}, nil
	}

	// Determine scope based on email template name
	// For now, we'll use ScopeAgent as default since most emails are agent-related
	// This can be enhanced to determine scope from email template metadata
	scope := constant.ScopeAgent

	// Update status to pending
	emailLog.StatusID = constant.StatusEmailPendingID
	if err := eu.emailRepo.UpdateStatusEmailLog(ctx, emailLog); err != nil {
		logger.Error(ctx, "failed to update email log status to pending", err.Error())
		return &emaildto.RetryEmailResponse{
			Success: false,
			Message: "Failed to update email log status",
		}, err
	}

	// Send email
	bodyText := "Please view this email in HTML format."
	err = eu.emailSender.Send(ctx, scope, emailLog.To, emailLog.Subject, emailLog.Body, bodyText)
	
	statusEmailID := constant.StatusEmailSuccessID
	if err != nil {
		logger.Error(ctx, "failed to retry email", err.Error())
		statusEmailID = constant.StatusEmailFailedID
		
		// Update metadata with error
		if emailLog.Meta == nil {
			emailLog.Meta = &entity.MetadataEmailLog{}
		}
		if emailLog.Meta.Notes != "" {
			emailLog.Meta.Notes += "; Retry failed: " + err.Error()
		} else {
			emailLog.Meta.Notes = "Retry failed: " + err.Error()
		}
	}

	// Update status
	emailLog.StatusID = uint(statusEmailID)
	if err := eu.emailRepo.UpdateStatusEmailLog(ctx, emailLog); err != nil {
		logger.Error(ctx, "failed to update email log status after retry", err.Error())
		return &emaildto.RetryEmailResponse{
			Success: false,
			Message: "Failed to update email log status",
		}, err
	}

	if err != nil {
		return &emaildto.RetryEmailResponse{
			Success: false,
			Message: "Failed to resend email: " + err.Error(),
		}, nil
	}

	return &emaildto.RetryEmailResponse{
		Success: true,
		Message: "Email resent successfully",
	}, nil
}

