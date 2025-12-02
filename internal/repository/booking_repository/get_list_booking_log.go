package booking_repository

import (
	"context"
	"encoding/json"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (br *BookingRepository) GetListBookingLog(ctx context.Context, filter *filter.BookingFilter) ([]entity.BookingDetail, int64, error) {
	db := br.db.GetTx(ctx)

	var bookingDetails []model.BookingDetail
	query := db.WithContext(ctx).Model(&model.BookingDetail{})

	if filter.Search != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("sub_booking_id ILIKE ? ", "%"+safeSearch+"%")
		query = query.Or("guest ILIKE ? ", "%"+safeSearch+"%")
		query = query.Or("detail_room->>'hotel_name' ILIKE ? ", "%"+safeSearch+"%")
	}

	if filter.BookingStatusID > 0 {
		query = query.Where("status_booking_id = ?", filter.BookingStatusID)
	}

	if filter.PaymentStatusID > 0 {
		query = query.Where("status_payment_id = ?", filter.PaymentStatusID)
	}

	if filter.ConfirmDateFrom != "" {
		query = query.Where("approved_at >= ?", filter.ConfirmDateFrom)
	}
	if filter.ConfirmDateTo != "" {
		query = query.Where("approved_at <= ?", filter.ConfirmDateTo)
	}

	if filter.CheckInDateFrom != "" {
		query = query.Where("check_in_date >= ?", filter.CheckInDateFrom)
	}
	if filter.CheckInDateTo != "" {
		query = query.Where("check_in_date <= ?", filter.CheckInDateTo)
	}
	if filter.CheckOutDateFrom != "" {
		query = query.Where("check_out_date >= ?", filter.CheckInDateFrom)
	}
	if filter.CheckOutDateTo != "" {
		query = query.Where("check_out_date <= ?", filter.CheckOutDateTo)
	}

	query = query.Where("approved_at IS NOT NULL")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting booking logs", err.Error())
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

	// Apply sorting
	if filter.Sort != "" {
		if validateColumnSort[filter.Sort] {
			query = query.Order(filter.Sort)
		}
	}
	query = query.Order("created_at DESC")

	if err := query.
		Preload("Booking").
		Preload("Booking.Agent").
		Preload("StatusBooking").
		Preload("StatusPayment").
		Find(&bookingDetails).Error; err != nil {
		logger.Error(ctx, "Error fetching booking logs", err.Error())
		return nil, 0, err
	}

	var result []entity.BookingDetail
	for _, detail := range bookingDetails {
		bookingLog := entity.BookingDetail{
			SubBookingID:  detail.SubBookingID,
			ApprovedAt:    detail.ApprovedAt,
			BookingStatus: detail.StatusBooking.Status,
			PaymentStatus: detail.StatusPayment.Status,
			CheckInDate:   detail.CheckInDate,
			CheckOutDate:  detail.CheckOutDate,
			Booking: entity.Booking{
				BookingCode: detail.Booking.BookingCode,
				AgentName:   detail.Booking.Agent.FullName,
			},
		}
		var detailRooms entity.DetailRoom
		if err := json.Unmarshal(detail.DetailRoom, &detailRooms); err != nil {
			logger.Error(ctx, "Error marshalling room detail to JSON", err.Error())
			detailRooms = entity.DetailRoom{}
		}
		bookingLog.DetailRooms = detailRooms
		result = append(result, bookingLog)
	}

	return result, total, nil
}
