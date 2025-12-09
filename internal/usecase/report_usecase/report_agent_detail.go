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

func (ru *ReportUsecase) ReportAgentDetail(ctx context.Context, req *reportdto.ReportAgentDetailRequest) (*reportdto.ReportAgentDetailResponse, error) {
	var dateFrom, dateTo *time.Time
	// Parse dates
	if req.DateFrom != "" {
		dateFromDt, err := time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from: %s", err.Error())
		}
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
		dateToDt = dateToDt.AddDate(0, 0, 1)
		dateTo = &dateToDt
	} else {
		now := time.Now().In(constant.AsiaJakarta)
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, constant.AsiaJakarta)
		dateTo = &endOfMonth
	}

	filterReq := filter.ReportDetailFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	filterReq.PaginationRequest = req.PaginationRequest
	if req.AgentID > 0 {
		filterReq.AgentID = &req.AgentID
	}
	if req.HotelID > 0 {
		filterReq.HotelID = &req.HotelID
	}

	data, total, err := ru.reportRepo.ReportAgentBookingDetail(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get report agent detail", err.Error())
		return nil, err
	}

	resp := &reportdto.ReportAgentDetailResponse{
		Total:                 total,
		ReportAgentDetailData: data,
	}

	return resp, nil
}
