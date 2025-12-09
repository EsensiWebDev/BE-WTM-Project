package report_repository

import (
	"context"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (rr *ReportRepository) ReportAgentBooking(ctx context.Context, filter filter.ReportFilter) ([]entity.ReportAgentBooking, int64, error) {
	db := rr.db.GetTx(ctx)

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Date filters with correct columns per status
	var dateConditions []string
	if filter.DateFrom != nil && filter.DateTo != nil {
		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 3 AND bd.approved_at >= $%d AND bd.approved_at < $%d)", argIndex, argIndex+1),
		)
		args = append(args, filter.DateFrom, filter.DateTo)
		argIndex += 2

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 5 AND bd.cancelled_at >= $%d AND bd.cancelled_at < $%d)", argIndex, argIndex+1),
		)
		args = append(args, filter.DateFrom, filter.DateTo)
		argIndex += 2

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 4 AND bd.rejected_at >= $%d AND bd.rejected_at < $%d)", argIndex, argIndex+1),
		)
		args = append(args, filter.DateFrom, filter.DateTo)
		argIndex += 2

		conditions = append(conditions, "("+strings.Join(dateConditions, " OR ")+")")
	} else if filter.DateFrom != nil {
		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 3 AND bd.approved_at >= $%d)", argIndex),
		)
		args = append(args, filter.DateFrom)
		argIndex++

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 5 AND bd.cancelled_at >= $%d)", argIndex),
		)
		args = append(args, filter.DateFrom)
		argIndex++

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 4 AND bd.rejected_at >= $%d)", argIndex),
		)
		args = append(args, filter.DateFrom)
		argIndex++

		conditions = append(conditions, "("+strings.Join(dateConditions, " OR ")+")")
	} else if filter.DateTo != nil {
		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 3 AND bd.approved_at < $%d)", argIndex),
		)
		args = append(args, filter.DateTo)
		argIndex++

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 5 AND bd.cancelled_at < $%d)", argIndex),
		)
		args = append(args, filter.DateTo)
		argIndex++

		dateConditions = append(dateConditions,
			fmt.Sprintf("(bd.status_booking_id = 4 AND bd.rejected_at < $%d)", argIndex),
		)
		args = append(args, filter.DateTo)
		argIndex++

		conditions = append(conditions, "("+strings.Join(dateConditions, " OR ")+")")
	}

	// Search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		conditions = append(conditions, fmt.Sprintf("u.full_name ILIKE $%d", argIndex))
		args = append(args, "%"+safeSearch+"%")
		argIndex++
	}

	// Hotel ID filter
	if len(filter.HotelID) > 0 {
		placeholders := make([]string, len(filter.HotelID))
		for i, id := range filter.HotelID {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, id)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("h.id IN (%s)", strings.Join(placeholders, ",")))
	}

	// Agent Company ID filter
	if len(filter.AgentCompanyID) > 0 {
		placeholders := make([]string, len(filter.AgentCompanyID))
		for i, id := range filter.AgentCompanyID {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, id)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("ac.id IN (%s)", strings.Join(placeholders, ",")))
	}

	// Build WHERE clause
	whereClause := "WHERE TRUE"
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Base query
	baseQuery := fmt.Sprintf(`
        SELECT
            h.id AS hotel_id,
            h.name AS hotel_name,
            ac.name AS agent_company,
            u.id AS agent_id,
            u.full_name AS agent_name,
            SUM(CASE WHEN bd.status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking,
            SUM(CASE WHEN bd.status_booking_id = 5 THEN 1 ELSE 0 END) AS cancelled_booking,
            SUM(CASE WHEN bd.status_booking_id = 4 THEN 1 ELSE 0 END) AS rejected_booking
        FROM booking_details bd
        JOIN bookings b ON bd.booking_id = b.id
        JOIN users u ON b.agent_id = u.id
        LEFT JOIN agent_companies ac ON u.agent_company_id = ac.id
        JOIN room_prices rp ON rp.id = bd.room_price_id
        JOIN room_types rt ON rp.room_type_id = rt.id
        JOIN hotels h ON rt.hotel_id = h.id
        %s
        GROUP BY h.id, h.name, ac.name, u.id, u.full_name
        ORDER BY h.id, u.id
    `, whereClause)

	// Count total records
	countQuery := fmt.Sprintf(`
        SELECT COUNT(*) FROM (
            %s
        ) subquery
    `, baseQuery)

	var total int64
	if err := db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		logger.Error(ctx, "Error counting report agent bookings", err.Error())
		return nil, 0, err
	}

	// Add pagination
	finalQuery := baseQuery
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		finalQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)
	}

	// Execute final query
	var reports []entity.ReportAgentBooking
	if err := db.WithContext(ctx).Raw(finalQuery, args...).Scan(&reports).Error; err != nil {
		logger.Error(ctx, "Error fetching report agent bookings", err.Error())
		return nil, 0, err
	}

	return reports, total, nil
}
