package report_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"

	"gorm.io/gorm/clause"
)

var validateColumnSort = map[string]bool{
	"status_booking_id": true,
	"check_in_date":     true,
	"check_out_date":    true,
	"guest":             true,
}

func (rr *ReportRepository) ReportAgentBookingDetail(ctx context.Context, filter filter.ReportDetailFilter) ([]entity.ReportAgentDetail, int64, error) {
	db := rr.db.GetTx(ctx)

	var bookingDetails []model.BookingDetail
	query := db.WithContext(ctx).
		Model(&model.BookingDetail{}).
		Preload("BookingDetailsAdditional").
		Preload("StatusBooking").
		Where("booking_details.status_booking_id IN ?", []int{constant.StatusBookingConfirmedID, constant.StatusBookingCancelledID, constant.StatusBookingRejectedID})

	if filter.HotelID != nil {
		query = query.
			Joins("JOIN room_prices rp ON booking_details.room_price_id = rp.id").
			Joins("JOIN room_types rt ON rp.room_type_id = rt.id").
			Where("rt.hotel_id = ?", *filter.HotelID)
	}
	if filter.AgentID != nil {
		query = query.Joins("JOIN bookings b ON booking_details.booking_id = b.id").
			Where("b.agent_id = ?", *filter.AgentID)
	}
	// Date filters with correct logic per status
	if filter.DateFrom != nil && filter.DateTo != nil {
		query = query.Where(
			db.Where("booking_details.status_booking_id = ? AND booking_details.approved_at >= ? AND booking_details.approved_at < ?",
				constant.StatusBookingConfirmedID, *filter.DateFrom, *filter.DateTo).
				Or("booking_details.status_booking_id = ? AND booking_details.cancelled_at >= ? AND booking_details.cancelled_at < ?",
					constant.StatusBookingCancelledID, *filter.DateFrom, *filter.DateTo).
				Or("booking_details.status_booking_id = ? AND booking_details.rejected_at >= ? AND booking_details.rejected_at < ?",
					constant.StatusBookingRejectedID, *filter.DateFrom, *filter.DateTo),
		)
	} else if filter.DateFrom != nil {
		query = query.Where(
			db.Where("booking_details.status_booking_id = ? AND booking_details.approved_at >= ?",
				constant.StatusBookingConfirmedID, *filter.DateFrom).
				Or("booking_details.status_booking_id = ? AND booking_details.cancelled_at >= ?",
					constant.StatusBookingCancelledID, *filter.DateFrom).
				Or("booking_details.status_booking_id = ? AND booking_details.rejected_at >= ?",
					constant.StatusBookingRejectedID, *filter.DateFrom),
		)
	} else if filter.DateTo != nil {
		query = query.Where(
			db.Where("booking_details.status_booking_id = ? AND booking_details.approved_at < ?",
				constant.StatusBookingConfirmedID, *filter.DateTo).
				Or("booking_details.status_booking_id = ? AND booking_details.cancelled_at < ?",
					constant.StatusBookingCancelledID, *filter.DateTo).
				Or("booking_details.status_booking_id = ? AND booking_details.rejected_at < ?",
					constant.StatusBookingRejectedID, *filter.DateTo),
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting booking details", err.Error())
		return nil, 0, err
	}

	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	if filter.Sort != "" {
		if validateColumnSort[filter.Sort] {
			var desc bool
			if strings.TrimSpace(strings.ToLower(filter.Dir)) == "asc" {
				desc = false
			} else {
				desc = true
			}
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: filter.Sort}, Desc: desc})
		}
	} else {
		query = query.Order("booking_details.id DESC")
	}

	if err := query.Find(&bookingDetails).Error; err != nil {
		logger.Error(ctx, "Error fetching booking details", err.Error())
		return nil, total, err
	}

	var results []entity.ReportAgentDetail
	for _, bd := range bookingDetails {
		var additionalNames []string
		for _, add := range bd.BookingDetailsAdditional {
			additionalNames = append(additionalNames, add.NameAdditional)
		}

		var detailRoom entity.DetailRoom
		if err := json.Unmarshal(bd.DetailRoom, &detailRoom); err != nil {
			logger.Error(ctx, "Error unmarshaling detail room", err.Error())
		}

		result := entity.ReportAgentDetail{
			GuestName:     bd.Guest,
			RoomType:      detailRoom.RoomTypeName,
			DateIn:        bd.CheckInDate.Format("2006-01-02"),
			DateOut:       bd.CheckOutDate.Format("2006-01-02"),
			Capacity:      fmt.Sprintf("%d Adult", detailRoom.Capacity),
			Additional:    strings.Join(additionalNames, ","),
			StatusBooking: bd.StatusBooking.Status,
		}
		results = append(results, result)
	}

	return results, total, nil

}
