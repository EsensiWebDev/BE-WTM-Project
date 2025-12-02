package report_repository

import (
	"context"
	"github.com/Masterminds/squirrel"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (rr *ReportRepository) ReportAgentBooking(ctx context.Context, filter filter.ReportFilter) ([]entity.ReportAgentBooking, int64, error) {
	db := rr.db.GetTx(ctx)

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar) // pakai $1, $2, ... untuk PostgreSQL

	// Base query builder
	builder := psql.Select(
		"h.id AS hotel_id",
		"h.name AS hotel_name",
		"ac.name AS agent_company",
		"u.id AS agent_id",
		"u.full_name AS agent_name",
		"SUM(CASE WHEN bd.status_booking_id = 3 THEN 1 END) AS confirmed_booking",
		"SUM(CASE WHEN bd.status_booking_id = 5 THEN 1 END) AS cancelled_booking",
	).From("booking_details bd").
		Join("bookings b ON bd.booking_id = b.id").
		Join("users u ON b.agent_id = u.id").
		LeftJoin("agent_companies ac ON u.agent_company_id = ac.id").
		Join("room_prices rp ON rp.id = bd.room_price_id").
		Join("room_types rt ON rp.room_type_id = rt.id").
		Join("hotels h ON rt.hotel_id = h.id").
		Where("TRUE").
		GroupBy("h.name, ac.name, u.full_name, h.id, u.id ")

	// Apply search filter
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		builder = builder.Where(
			squirrel.Or{
				squirrel.Expr("u.full_name ILIKE ? ", "%"+safeSearch+"%"),
			},
		)
	}

	// Apply other filters
	if filter.DateFrom != nil {
		builder = builder.Where(squirrel.GtOrEq{"b.approved_at": filter.DateFrom.Format("2006-01-02")})
	}
	if filter.DateTo != nil {
		builder = builder.Where(squirrel.LtOrEq{"b.approved_at": filter.DateTo.Format("2006-01-02")})
	}

	if len(filter.HotelID) > 0 {
		builder = builder.Where(squirrel.Eq{"h.id": filter.HotelID})
	}

	if len(filter.AgentCompanyID) > 0 {
		builder = builder.Where(squirrel.Eq{"ac.id": filter.AgentCompanyID})
	}

	// Count total records
	var total int64
	countBuilder := squirrel.Select("COUNT(*)").FromSelect(builder, "subquery")
	countQuery, countArgs, err := countBuilder.ToSql()
	if err := db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
		logger.Error(ctx, "Error counting report agent bookings", err.Error())
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		builder = builder.Limit(uint64(filter.Limit)).Offset(uint64(offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	// Execute final query
	var reports []entity.ReportAgentBooking
	if err := db.WithContext(ctx).Raw(query, args...).Debug().Scan(&reports).Error; err != nil {
		logger.Error(ctx, "Error fetching report agent bookings", err.Error())
		return nil, 0, err
	}

	return reports, total, nil
}
