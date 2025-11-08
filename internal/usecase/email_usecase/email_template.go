package email_usecase

import (
	"context"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) EmailTemplate(ctx context.Context) (*emaildto.EmailTemplateResponse, error) {

	emailTemplate, err := eu.emailRepo.GetEmailTemplateByName(ctx, constant.EmailHotelBookingRequest)
	if err != nil {
		logger.Error(ctx, "Error getting email template by name:", err.Error())
		return nil, err
	}

	if emailTemplate == nil {
		logger.Error(ctx, "Email template not found for name:", constant.EmailHotelBookingRequest)
		return nil, nil
	}

	resp := &emaildto.EmailTemplateResponse{
		Body:      emailTemplate.Body,
		Subject:   emailTemplate.Subject,
		Signature: emailTemplate.Signature,
	}

	if emailTemplate.IsSignatureImage {

	}

	return resp, nil
}
