package report_repository

import (
	"context"
	"encoding/json"
	"gorm.io/gorm/clause"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
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
		Preload("StatusBooking")

	if filter.HotelID != nil {
		query = query.Joins("JOIN room_types rt ON booking_details.room_type_id = rt.id").
			Where("rt.hotel_id = ?", *filter.HotelID)
	}
	if filter.AgentID != nil {
		query = query.Joins("JOIN bookings b ON booking_details.booking_id = b.id").
			Where("b.agent_id = ?", *filter.AgentID)
	}

	var total int64
	if err := query.Count(&total).Debug().Error; err != nil {
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

	if err := query.Find(&bookingDetails).Debug().Error; err != nil {
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
			logger.Error(ctx, "Error marshalling room detail to JSON", err.Error())
			detailRoom = entity.DetailRoom{}
		}

		result := entity.ReportAgentDetail{
			GuestName:     bd.Guest,
			RoomType:      detailRoom.RoomTypeName,
			DateIn:        bd.CheckInDate.Format("2006-01-02"),
			DateOut:       bd.CheckOutDate.Format("2006-01-02"),
			Capacity:      detailRoom.Capacity,
			Additional:    strings.Join(additionalNames, ","),
			StatusBooking: bd.StatusBooking.Status,
		}
		results = append(results, result)
	}

	return results, total, nil

}
