package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookingLog(ctx context.Context, req *bookingdto.ListBookingLogRequest) (*bookingdto.ListBookingLogResponse, error) {
	filterReq := filter.BookingFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	filterReq.ConfirmDateFrom = req.ConfirmDateFrom
	filterReq.ConfirmDateTo = req.ConfirmDateTo

	filterReq.CheckInDateFrom = req.CheckInDateFrom
	filterReq.CheckInDateTo = req.CheckInDateTo

	filterReq.CheckOutDateFrom = req.CheckOutDateFrom
	filterReq.CheckOutDateTo = req.CheckOutDateTo

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
