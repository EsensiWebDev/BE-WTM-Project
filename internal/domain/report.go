package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/reportdto"
	"wtm-backend/internal/repository/filter"
)

type ReportUsecase interface {
	ReportAgent(ctx context.Context, req *reportdto.ReportRequest) (*reportdto.ReportAgentResponse, error)
	ReportAgentDetail(ctx context.Context, req *reportdto.ReportAgentDetailRequest) (*reportdto.ReportAgentDetailResponse, error)
	ReportSummary(ctx context.Context, req *reportdto.ReportSummaryRequest) (*reportdto.ReportSummaryResponse, error)
}

type ReportRepository interface {
	ReportAgentBooking(ctx context.Context, filter filter.ReportFilter) ([]entity.ReportAgentBooking, int64, error)
	ReportAgentBookingDetail(ctx context.Context, filter filter.ReportDetailFilter) ([]entity.ReportAgentDetail, int64, error)
	ReportBookingSummary(ctx context.Context, filter filter.ReportSummaryFilter) ([]entity.MonthlyBookingSummary, error)
	ReportForGraph(ctx context.Context, filter filter.ReportSummaryFilter) ([]entity.ReportForGraph, error)
}
