package email_usecase

import (
	"context"
	"time"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) GetEmailLogDetail(ctx context.Context, id uint) (*emaildto.EmailLogDetailResponse, error) {
	emailLog, err := eu.emailRepo.GetEmailLogByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed to get email log detail", err.Error())
		return nil, err
	}

	response := &emaildto.EmailLogDetailResponse{
		ID:        emailLog.ID,
		To:        emailLog.To,
		Subject:   emailLog.Subject,
		Body:      emailLog.Body,
		DateTime:  emailLog.CreatedAt.Format(time.RFC3339),
		EmailType: emailLog.EmailType,
		Status:    constant.MapStatusEmailLog[int(emailLog.StatusID)],
	}

	if emailLog.Meta != nil {
		response.HotelName = emailLog.Meta.HotelName
		response.Notes = emailLog.Meta.Notes
	}

	return response, nil
}

