package report_usecase

import "wtm-backend/internal/domain"

type ReportUsecase struct {
	reportRepo domain.ReportRepository
}

func NewReportUsecase(reportRepo domain.ReportRepository) *ReportUsecase {
	return &ReportUsecase{
		reportRepo: reportRepo,
	}
}
