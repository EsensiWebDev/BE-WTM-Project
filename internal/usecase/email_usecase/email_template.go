package email_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) EmailTemplate(ctx context.Context, req *emaildto.EmailTemplateRequest) (*emaildto.EmailTemplateResponse, error) {

	var nameTemplate string
	switch req.Type {
	case "confirm":
		nameTemplate = constant.EmailHotelBookingRequest
	case "cancel":
		nameTemplate = constant.EmailHotelBookingCancel
	default:
		nameTemplate = constant.EmailHotelBookingRequest
	}

	emailTemplate, err := eu.emailRepo.GetEmailTemplateByName(ctx, nameTemplate)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return nil, err
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for name:", nameTemplate)
		return nil, nil
	}

	resp := &emaildto.EmailTemplateResponse{
		Body:      emailTemplate.Body,
		Subject:   emailTemplate.Subject,
		Signature: emailTemplate.Signature,
	}

	if emailTemplate.IsSignatureImage {
		bucketName := fmt.Sprintf("%s-%s", constant.ConstEmail, constant.ConstPublic)
		signatureImg, err := eu.fileStorage.GetFile(ctx, bucketName, emailTemplate.Signature)
		if err != nil {
			logger.Error(ctx, "Error getting signature image:", err.Error())
		}
		resp.Signature = signatureImg
	}

	return resp, nil
}
