package booking_usecase

import (
	"context"
	"fmt"
	"time"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookingLog(ctx context.Context, req *bookingdto.ListBookingLogRequest) (*bookingdto.ListBookingLogResponse, error) {
	filterReq := filter.BookingFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	//validate date
	if req.ConfirmDateFrom != "" {
		confirmDateFrom, err := time.Parse(time.DateOnly, req.ConfirmDateFrom)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !confirmDateFrom.IsZero() {
			confirmDateFrom = confirmDateFrom.Truncate(time.Hour * 24)
			filterReq.ConfirmDateFrom = &confirmDateFrom
		}
	}

	if req.ConfirmDateTo != "" {
		confirmDateTo, err := time.Parse(time.DateOnly, req.ConfirmDateTo)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !confirmDateTo.IsZero() {
			confirmDateTo = confirmDateTo.Truncate(time.Hour*24).AddDate(0, 0, 1)
			filterReq.ConfirmDateTo = &confirmDateTo
		}
	}

	if req.CheckInDateFrom != "" {
		checkInDateFrom, err := time.Parse(time.DateOnly, req.CheckInDateFrom)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !checkInDateFrom.IsZero() {
			checkInDateFrom = checkInDateFrom.Truncate(time.Hour * 24)
			filterReq.CheckInDateFrom = &checkInDateFrom
		}
	}

	if req.CheckInDateTo != "" {
		checkInDateTo, err := time.Parse(time.DateOnly, req.CheckInDateTo)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !checkInDateTo.IsZero() {
			checkInDateTo = checkInDateTo.Truncate(time.Hour*24).AddDate(0, 0, 1)
			filterReq.CheckInDateTo = &checkInDateTo
		}
	}

	if req.CheckOutDateFrom != "" {
		checkOutDateFrom, err := time.Parse(time.DateOnly, req.CheckOutDateFrom)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !checkOutDateFrom.IsZero() {
			checkOutDateFrom = checkOutDateFrom.Truncate(time.Hour * 24)
			filterReq.CheckInDateFrom = &checkOutDateFrom
		}
	}

	if req.CheckOutDateTo != "" {
		checkOutDateTo, err := time.Parse(time.DateOnly, req.CheckOutDateTo)
		if err != nil {
			logger.Error(ctx, "ListHotelsForAgent", err.Error())
		}
		if !checkOutDateTo.IsZero() {
			checkOutDateTo = checkOutDateTo.Truncate(time.Hour*24).AddDate(0, 0, 1)
			filterReq.CheckInDateTo = &checkOutDateTo
		}

	}

	filterReq.BookingStatusID = req.BookingStatusID
	filterReq.PaymentStatusID = req.PaymentStatusID

	subBookings, total, err := bu.bookingRepo.GetListBookingLog(ctx, &filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get bookings", err.Error())
		return nil, err
	}

	resp := &bookingdto.ListBookingLogResponse{
		Total: total,
		Data:  make([]bookingdto.BookingLog, 0),
	}

	for _, detail := range subBookings {
		bookingLog := bookingdto.BookingLog{
			SubBookingID:  detail.SubBookingID,
			BookingID:     detail.Booking.BookingCode,
			AgentName:     detail.Booking.AgentName,
			BookingStatus: detail.BookingStatus,
			PaymentStatus: detail.PaymentStatus,
			CheckInDate:   detail.CheckInDate.Format("2006-01-02"),
			CheckOutDate:  detail.CheckOutDate.Format("2006-01-02"),
			HotelName:     detail.DetailRooms.HotelName,
			RoomTypeName:  detail.DetailRooms.RoomTypeName,
			RoomNights:    int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24),
			Capacity:      fmt.Sprintf("%d Adult", detail.DetailRooms.Capacity),
			ConfirmDate:   detail.ApprovedAt.Format("2006-01-02"),
		}
		if detail.ApprovedAt.IsZero() {
			bookingLog.ConfirmDate = "N/A"
		}
		resp.Data = append(resp.Data, bookingLog)
	}

	return resp, nil
}
