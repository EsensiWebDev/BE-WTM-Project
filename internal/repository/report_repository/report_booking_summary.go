package report_repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (rr *ReportRepository) ReportBookingSummary(ctx context.Context, filter filter.ReportFilter) ([]entity.MonthlyBookingSummary, error) {
	db := rr.db.GetTx(ctx)
	var summaries []entity.MonthlyBookingSummary

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	if !filter.IsRangeDate {
		builder := psql.Select(
			"DATE_TRUNC('month', bd.approved_at) AS month",
			"SUM(CASE WHEN bd.status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking",
			"SUM(CASE WHEN bd.status_booking_id = 5 THEN 1 ELSE 0 END) AS cancellation_booking",
		).From("booking_details bd").
			Where(squirrel.GtOrEq{"bd.approved_at": filter.DateFrom}).
			Where(squirrel.Lt{"bd.approved_at": filter.DateTo})

		builder = builder.GroupBy("DATE_TRUNC('month', bd.approved_at)").OrderBy("month")

		query, args, err := builder.ToSql()
		if err != nil {
			logger.Error(ctx, "Error building booking summary query", err.Error())
			return nil, err
		}

		if err := db.WithContext(ctx).Raw(query, args...).Scan(&summaries).Debug().Error; err != nil {
			logger.Error(ctx, "Error executing booking summary query", err.Error())
			return nil, err
		}

		return summaries, nil

	}

	// This Month Summary
	builderThisMonth := psql.Select(
		"'This Month' AS month",
		"SUM(CASE WHEN bd.status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking",
		"SUM(CASE WHEN bd.status_booking_id = 5 THEN 1 ELSE 0 END) AS cancellation_booking",
	).From("booking_details bd").
		Where(squirrel.GtOrEq{"bd.approved_at": filter.DateFrom}).
		Where(squirrel.Lt{"bd.approved_at": filter.DateTo})

	queryThisMonth, argsThisMonth, err := builderThisMonth.ToSql()
	if err != nil {
		logger.Error(ctx, "Error building booking summary for this month query", err.Error())
		return nil, err
	}

	var summaryThisMonth entity.MonthlyBookingSummary
	if err := db.WithContext(ctx).Raw(queryThisMonth, argsThisMonth...).Scan(&summaryThisMonth).Debug().Error; err != nil {
		logger.Error(ctx, "Error executing booking summary for this month query", err.Error())
		return nil, err
	}

	// Total Summary
	builderTotal := psql.Select(
		"'Total' AS month",
		"SUM(CASE WHEN bd.status_booking_id = 3 THEN 1 ELSE 0 END) AS confirmed_booking",
		"SUM(CASE WHEN bd.status_booking_id = 5 THEN 1 ELSE 0 END) AS cancellation_booking",
	).From("booking_details bd").
		Where("TRUE")

	queryTotal, argsTotal, err := builderTotal.ToSql()
	if err != nil {
		logger.Error(ctx, "Error building total booking summary query", err.Error())
		return nil, err
	}

	var summaryTotal entity.MonthlyBookingSummary
	if err := db.WithContext(ctx).Raw(queryTotal, argsTotal...).Debug().Scan(&summaryTotal).Debug().Error; err != nil {
		logger.Error(ctx, "Error executing total booking summary query", err.Error())
		return nil, err
	}

	summaries = append(summaries, summaryThisMonth, summaryTotal)

	return summaries, nil
}
