package user_usecase

import (
	"context"
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
	var loginLink = "https://hotelbox.com/login"
	var reRegisterLink = "https://hotelbox.com/re-register"

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

	err = uu.emailSender.Send(ctx, email, emailTemplate.Subject, bodyHTML, bodyText)
	if err != nil {
		logger.Error(ctx, "Error sending email:", err.Error())
	}

}

type EmailData struct {
	AgentName       string
	LoginLink       string
	ReRegisterLink  string
	RejectionReason string
}
