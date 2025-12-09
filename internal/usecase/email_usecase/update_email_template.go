package email_usecase

import (
	"context"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) UpdateEmailTemplate(ctx context.Context, req *emaildto.UpdateEmailTemplateRequest) error {

	var typeName string
	switch req.Type {
	case "confirm":
		typeName = constant.EmailHotelBookingRequest
	case "cancel":
		typeName = constant.EmailHotelBookingCancel
	default:
		typeName = constant.EmailHotelBookingRequest
	}

	emailTemplate, err := eu.emailRepo.GetEmailTemplateByName(ctx, typeName)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return err
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for name:", constant.EmailHotelBookingRequest)
		return nil
	}

	if req.Subject != "" {
		emailTemplate.Subject = req.Subject
	}
	if req.Body != "" {
		emailTemplate.Body = req.Body
	}
	if req.SignatureImage != nil {
		emailTemplate.IsSignatureImage = true
		url, err := eu.uploadFile(ctx, req.SignatureImage)
		if err != nil {
			logger.Error(ctx, "Error uploading signature image:", err.Error())
			return err
		}
		emailTemplate.Signature = url
	} else if req.SignatureText != "" {
		emailTemplate.Signature = req.SignatureText
	}

	return eu.emailRepo.UpdateEmailTemplate(ctx, emailTemplate)
}
