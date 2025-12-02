package booking_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (br *BookingRepository) GetBookingIDs(ctx context.Context, agentID uint, filter *filter.DefaultFilter) ([]string, int64, error) {
	db := br.db.GetTx(ctx)

	query := db.WithContext(ctx).
		Model(&model.Booking{}).
		Select("booking_code").
		Where("agent_id = ?", agentID)

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting bed types", err.Error())
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// apply sort
	query = query.Order("created_at desc")

	var bookingIDs []string
	if err := query.Pluck("booking_code", &bookingIDs).Error; err != nil {
		logger.Error(ctx, "Error fetching bed types", err.Error())
		return nil, total, err
	}

	return bookingIDs, total, nil

}
