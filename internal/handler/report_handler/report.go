package report_handler

import (
	"wtm-backend/internal/domain"
)

type ReportHandler struct {
	reportUsecase domain.ReportUsecase
}

func NewReportHandler(reportUsecase domain.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
	}
}
