package user_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/userdto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (uu *UserUsecase) UpdateStatusUser(ctx context.Context, req *userdto.UpdateStatusUserRequest) error {
	var statusUser uint
	if req.IsActive {
		statusUser = constant.StatusUserActiveID
	} else {
		statusUser = constant.StatusUserRejectID
	}

	user, err := uu.userRepo.UpdateStatusUser(ctx, req.ID, statusUser)
	if err != nil {
		logger.Error(ctx, "Error updating user status:", err.Error())
		return err
	}

	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), uu.config.DurationCtxTOSlow)
		defer cancel()
		uu.sendEmailNotification(newCtx, req, user.FullName, user.Email)
	}()

	return nil
}

func (uu *UserUsecase) sendEmailNotification(ctx context.Context, req *userdto.UpdateStatusUserRequest, name, email string) {
	var statusEmail string
	var loginLink = fmt.Sprintf("%s/login", uu.config.URLFEAgent)
	var reRegisterLink = fmt.Sprintf("%s/register", uu.config.URLFEAgent)

	if req.IsActive {
		statusEmail = constant.EmailAgentApproved
	} else {
		statusEmail = constant.EmailAgentRejected
	}

	emailTemplate, err := uu.emailRepo.GetEmailTemplateByName(ctx, statusEmail)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for status:", statusEmail)
		return
	}

	// Inject data
	data := EmailData{
		AgentName:       name,
		LoginLink:       loginLink,
		ReRegisterLink:  reRegisterLink,
		RejectionReason: req.Reason,
	}

	bodyHTML, err := utils.ParseTemplate(emailTemplate.Body, data)
	if err != nil {
		logger.Error(ctx, "Error parsing body HTML:", err.Error())
		return
	}

	bodyText := "Please view this email in HTML format." // Optional fallback

	subjectParsed := emailTemplate.Subject

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
	if err = uu.emailRepo.CreateEmailLog(ctx, &emailLog); err != nil {
		logger.Error(ctx, "Failed to create email log:", err)
		dataEmail = false
	} else {
		dataEmail = true
	}

	err = uu.emailSender.Send(ctx, constant.ScopeAgent, emailTo, subjectParsed, bodyHTML, bodyText)
	if err != nil {
		logger.Error(ctx, "Failed to sending email:", err.Error())
		statusEmailID = constant.StatusEmailFailedID
		metadataLog.Notes = fmt.Sprintf("Failed to sending email: %s", err.Error())
		emailLog.Meta = &metadataLog
	}

	if dataEmail {
		emailLog.StatusID = uint(statusEmailID)
		if err := uu.emailRepo.UpdateStatusEmailLog(ctx, &emailLog); err != nil {
			logger.Error(ctx, "Failed to update email log:", err.Error())
		}
	}

}

type EmailData struct {
	AgentName       string
	LoginLink       string
	ReRegisterLink  string
	RejectionReason string
}
