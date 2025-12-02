package report_repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (rr *ReportRepository) ReportForGraph(ctx context.Context, filter filter.ReportFilter) ([]entity.ReportForGraph, error) {
	db := rr.db.GetTx(ctx)
	var reports []entity.ReportForGraph

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	builder := psql.Select(
		"bd.approved_at AS date",
		"COUNT(bd.id) AS count",
	).From("booking_details bd").
		Where("TRUE")

	// Apply filters
	if filter.DateFrom != nil {
		builder = builder.Where(squirrel.GtOrEq{"bd.approved_at": filter.DateFrom})
	}
	if filter.DateTo != nil {
		builder = builder.Where(squirrel.LtOrEq{"bd.approved_at": filter.DateTo})
	}

	builder = builder.GroupBy("bd.approved_at").OrderBy("bd.approved_at")

	// Build query
	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(ctx, "Error building report for graph query", err.Error())
		return nil, err
	}

	// Execute query
	if err := db.WithContext(ctx).Raw(query, args...).Debug().Scan(&reports).Error; err != nil {
		logger.Error(ctx, "Error executing report for graph query", err.Error())
		return nil, err
	}

	return reports, nil
}
