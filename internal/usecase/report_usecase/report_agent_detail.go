package report_usecase

import (
	"context"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (ru *ReportUsecase) ReportAgentDetail(ctx context.Context, req *reportdto.ReportAgentDetailRequest) (*reportdto.ReportAgentDetailResponse, error) {
	filterReq := filter.ReportDetailFilter{}
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
