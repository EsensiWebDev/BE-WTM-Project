package booking_usecase

import (
	"context"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookingLog(ctx context.Context, req *bookingdto.ListBookingLogRequest) (*bookingdto.ListBookingLogResponse, error) {
	filterReq := filter.BookingFilter{}
	filterReq.PaginationRequest = req.PaginationRequest

	bookings, total, err := bu.bookingRepo.GetBookings(ctx, &filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get bookings", err.Error())
		return nil, err
	}

	resp := &bookingdto.ListBookingLogResponse{
		Total: total,
		Data:  make([]bookingdto.BookingLog, 0),
	}

	for _, booking := range bookings {
		for _, detail := range booking.BookingDetails {
			bookingLog := bookingdto.BookingLog{
				BookingCode:   booking.BookingCode,
				ConfirmDate:   detail.UpdatedAt.Format("2006-01-02"),
				AgentName:     booking.AgentName,
				BookingStatus: booking.BookingStatus,
				PaymentStatus: booking.PaymentStatus,
				CheckInDate:   detail.CheckInDate.Format("2006-01-02"),
				CheckOutDate:  detail.CheckOutDate.Format("2006-01-02"),
				HotelName:     detail.DetailRooms.HotelName,
				RoomTypeName:  detail.DetailRooms.RoomTypeName,
				RoomNights:    int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24),
				Capacity:      detail.DetailRooms.Capacity,
			}
			resp.Data = append(resp.Data, bookingLog)
		}
	}

	return resp, nil
}
