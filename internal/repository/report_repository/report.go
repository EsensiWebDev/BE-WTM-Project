package report_repository

import "wtm-backend/internal/infrastructure/database"

type ReportRepository struct {
	db *database.DBPostgre
}

func NewReportRepository(db *database.DBPostgre) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}
