package report_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (rr *ReportRepository) ReportBookingSummary(ctx context.Context, filter filter.ReportSummaryFilter) ([]entity.MonthlyBookingSummary, error) {
	db := rr.db.GetTx(ctx)
	var summaries []entity.MonthlyBookingSummary

	// This Month Summary
	queryThisMonth := `
        SELECT 
            'This Month' AS month,
            SUM(CASE WHEN status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking,
            SUM(CASE WHEN status_booking_id = 5 THEN 1 ELSE 0 END) AS cancelled_booking,
            SUM(CASE WHEN status_booking_id = 4 THEN 1 ELSE 0 END) AS rejected_booking
        FROM booking_details
        WHERE 
            (status_booking_id = 3 AND approved_at >= ? AND approved_at < ?)
            OR (status_booking_id = 5 AND cancelled_at >= ? AND cancelled_at < ?)
            OR (status_booking_id = 4 AND rejected_at >= ? AND rejected_at < ?)
    `

	var summaryThisMonth entity.MonthlyBookingSummary
	if err := db.WithContext(ctx).Raw(queryThisMonth,
		filter.DateFrom, filter.DateTo,
		filter.DateFrom, filter.DateTo,
		filter.DateFrom, filter.DateTo,
	).Scan(&summaryThisMonth).Error; err != nil {
		logger.Error(ctx, "Error executing booking summary for this month query", err.Error())
		return nil, err
	}

	// Total Summary (all time)
	queryTotal := `
        SELECT 
            'Total' AS month,
            SUM(CASE WHEN status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking,
            SUM(CASE WHEN status_booking_id = 5 THEN 1 ELSE 0 END) AS cancelled_booking,
            SUM(CASE WHEN status_booking_id = 4 THEN 1 ELSE 0 END) AS rejected_booking
        FROM booking_details
        WHERE 
            (status_booking_id = 3 AND approved_at IS NOT NULL)
            OR (status_booking_id = 5 AND cancelled_at IS NOT NULL)
            OR (status_booking_id = 4 AND rejected_at IS NOT NULL)
    `

	var summaryTotal entity.MonthlyBookingSummary
	if err := db.WithContext(ctx).Raw(queryTotal).Scan(&summaryTotal).Error; err != nil {
		logger.Error(ctx, "Error executing total booking summary query", err.Error())
		return nil, err
	}

	summaries = append(summaries, summaryThisMonth, summaryTotal)

	return summaries, nil
}
