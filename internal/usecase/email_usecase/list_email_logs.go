package email_usecase

import (
	"context"
	"time"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (eu *EmailUsecase) ListEmailLogs(ctx context.Context, req *emaildto.ListEmailLogsRequest) (*emaildto.ListEmailLogsResponse, error) {
	filterReq := filter.DefaultFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	emailLogs, total, err := eu.emailRepo.GetEmailLogs(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "failed to list email logs", err.Error())
		return nil, err
	}

	response := &emaildto.ListEmailLogsResponse{
		EmailLogs: make([]emaildto.EmailLogResponse, 0, len(emailLogs)),
		Total:     total,
	}

	for _, log := range emailLogs {
		logData := emaildto.EmailLogResponse{
			DateTime:  log.CreatedAt.Format(time.RFC3339),
			HotelName: log.Meta.HotelName,
			Status:    constant.MapStatusEmailLog[int(log.StatusID)],
			Notes:     log.Meta.Notes,
		}
		response.EmailLogs = append(response.EmailLogs, logData)
	}

	return response, nil
}
