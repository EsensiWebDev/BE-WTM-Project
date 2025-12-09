package report_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (rr *ReportRepository) ReportForGraph(ctx context.Context, filter filter.ReportSummaryFilter) ([]entity.ReportForGraph, error) {
	db := rr.db.GetTx(ctx)
	var reports []entity.ReportForGraph

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Only confirmed bookings
	conditions = append(conditions, fmt.Sprintf("bd.status_booking_id = %d", constant.StatusBookingConfirmedID))

	// Date filters
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("bd.approved_at >= $%d", argIndex))
		args = append(args, filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("bd.approved_at < $%d", argIndex))
		args = append(args, filter.DateTo)
		argIndex++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Build query
	query := fmt.Sprintf(`
        SELECT
            DATE(bd.approved_at) AS date_time,
            COUNT(bd.id) AS count
        FROM booking_details bd
        %s
        GROUP BY date_time
        ORDER BY date_time
    `, whereClause)

	// Execute query
	if err := db.WithContext(ctx).Raw(query, args...).Scan(&reports).Error; err != nil {
		logger.Error(ctx, "Error executing report for graph query", err.Error())
		return nil, err
	}

	return reports, nil
}
