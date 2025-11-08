package report_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (ru *ReportUsecase) ReportAgent(ctx context.Context, req *reportdto.ReportRequest) (*reportdto.ReportAgentResponse, error) {
	var dateFrom, dateTo *time.Time
	var isDateFromSet, isDateToSet bool
	// Parse dates
	if req.DateFrom != "" {
		dateFromDt, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from: %s", err.Error())
		}
		isDateFromSet = true
		dateFrom = &dateFromDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, constant.AsiaJakarta)
		dateFrom = &startOfMonth
	}

	if req.DateTo != "" {
		dateToDt, err := time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to: %s", err.Error())
		}
		isDateToSet = true
		dateTo = &dateToDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, constant.AsiaJakarta)
		dateTo = &endOfMonth
	}

	filterReq := filter.ReportFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	if req.HotelID > 0 {
		filterReq.HotelID = &req.HotelID
	}
	if req.AgentCompanyID > 0 {
		filterReq.AgentCompanyID = &req.AgentCompanyID
	}

	if isDateFromSet || isDateToSet {
		filterReq.IsRangeDate = true
	}

	filterReq.PaginationRequest = req.PaginationRequest
	reportAgent, total, err := ru.reportRepo.ReportAgentBooking(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get report agent", err.Error())
		return nil, err
	}
	resp := &reportdto.ReportAgentResponse{
		Total: total,
		Data:  reportAgent,
	}

	return resp, nil
}
