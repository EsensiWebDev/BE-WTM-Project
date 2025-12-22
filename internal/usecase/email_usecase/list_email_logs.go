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
	filterReq := filter.EmailLogFilter{}
	filterReq.PaginationRequest = req.PaginationRequest
	filterReq.EmailType = []string{constant.EmailHotelBookingRequest, constant.EmailHotelBookingCancel}

	// Apply filters from request
	if len(req.Status) > 0 {
		filterReq.Status = req.Status
	}
	if req.HotelName != "" {
		filterReq.HotelName = req.HotelName
	}
	if req.BookingCode != "" {
		filterReq.BookingCode = req.BookingCode
	}

	// Parse date strings
	if req.DateFrom != "" {
		if dateFrom, err := time.Parse("2006-01-02", req.DateFrom); err == nil {
			filterReq.DateFrom = &dateFrom
		} else if dateFrom, err := time.Parse(time.RFC3339, req.DateFrom); err == nil {
			filterReq.DateFrom = &dateFrom
		}
	}
	if req.DateTo != "" {
		if dateTo, err := time.Parse("2006-01-02", req.DateTo); err == nil {
			// Set to end of day
			dateTo = dateTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filterReq.DateTo = &dateTo
		} else if dateTo, err := time.Parse(time.RFC3339, req.DateTo); err == nil {
			filterReq.DateTo = &dateTo
		}
	}

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
			ID:        log.ID,
			DateTime:  log.CreatedAt.Format(time.RFC3339),
			Status:    constant.MapStatusEmailLog[int(log.StatusID)],
			EmailType: log.EmailType,
		}
		if log.Meta != nil {
			logData.HotelName = log.Meta.HotelName
			logData.Notes = log.Meta.Notes
			logData.BookingCode = log.Meta.BookingCode
		}
		response.EmailLogs = append(response.EmailLogs, logData)
	}

	return response, nil
}
