package report_repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
)

func (rr *ReportRepository) ReportNewAgentSummary(ctx context.Context, filter filter.ReportFilter) ([]entity.MonthlyNewAgentSummary, error) {
	db := rr.db.GetTx(ctx)
	var summaries []entity.MonthlyNewAgentSummary

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	if !filter.IsRangeDate {
		builder := psql.Select(
			"DATE_TRUNC('month', u.created_at) AS month",
			"COUNT(u.id) AS new_agent_count",
		).From("users u").
			Where(squirrel.Eq{"u.role_id": constant.RoleAgentID}).
			Where(squirrel.GtOrEq{"u.created_at": filter.DateFrom}).
			Where(squirrel.Lt{"u.created_at": filter.DateTo})

		builder = builder.GroupBy("DATE_TRUNC('month', u.created_at)").OrderBy("month")

		query, args, err := builder.ToSql()
		if err != nil {
			return nil, err
		}

		if err := db.WithContext(ctx).Raw(query, args...).Scan(&summaries).Error; err != nil {
			return nil, err
		}

		return summaries, nil

	}
	// This Month Summary
	builderThisMonth := psql.Select(
		"'This Month' AS month",
		"COUNT(u.id) AS new_agent_count",
	).From("users u").
		Where(squirrel.Eq{"u.role_id": constant.RoleAgentID}).
		Where(squirrel.GtOrEq{"u.created_at": filter.DateFrom}).
		Where(squirrel.Lt{"u.created_at": filter.DateTo})

	queryThisMonth, argsThisMonth, err := builderThisMonth.ToSql()
	if err != nil {
		return nil, err
	}

	var summaryThisMonth entity.MonthlyNewAgentSummary
	if err := db.WithContext(ctx).Raw(queryThisMonth, argsThisMonth...).Scan(&summaryThisMonth).Error; err != nil {
		return nil, err
	}

	// Total Summary
	builderTotal := psql.Select(
		"'Total' AS month",
		"COUNT(u.id) AS new_agent_count",
	).From("users u").
		Where(squirrel.Eq{"u.role_id": constant.RoleAgentID})

	queryTotal, argsTotal, err := builderTotal.ToSql()
	if err != nil {
		return nil, err
	}

	var summaryTotal entity.MonthlyNewAgentSummary
	if err := db.WithContext(ctx).Raw(queryTotal, argsTotal...).Debug().Scan(&summaryTotal).Error; err != nil {
		return nil, err
	}

	summaries = append(summaries, summaryThisMonth, summaryTotal)

	return summaries, nil
}
