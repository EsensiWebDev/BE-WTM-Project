package booking_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookingHistory(ctx context.Context, req *bookingdto.ListBookingHistoryRequest) (resp *bookingdto.ListBookingHistoryResponse, err error) {
	// Get agent Id from context
	userCtx, err := bu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return nil, fmt.Errorf("failed to get user from context: %s", err.Error())
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return nil, fmt.Errorf("user context is nil")
	}

	agentID := userCtx.ID

	bookingFilter := filter.BookingFilter{
		AgentID:         agentID,
		StatusBookingID: req.StatusBookingID,
		StatusPaymentID: req.StatusPaymentID,
	}
	if req.SearchBy == "booking_id" {
		bookingFilter.BookingIDSearch = req.Search
	} else if req.SearchBy == "guest_name" {
		bookingFilter.GuestNameSearch = req.Search
	}

	bookings, total, err := bu.bookingRepo.GetBookings(ctx, &bookingFilter)
	if err != nil {
		logger.Error(ctx, "failed to get bookings", err.Error())
		return nil, fmt.Errorf("failed to get bookings: %s", err.Error())
	}

	resp = &bookingdto.ListBookingHistoryResponse{
		Total: total,
		Data:  make([]bookingdto.DataBookingHistory, len(bookings)),
	}

	for i, booking := range bookings {
		resp.Data[i] = bookingdto.DataBookingHistory{
			BookingID:     booking.ID,
			GuestName:     booking.Guests,
			BookingCode:   booking.BookingCode,
			BookingStatus: booking.BookingStatus,
			PaymentStatus: booking.PaymentStatus,
			Detail:        make([]bookingdto.DetailBookingHistory, len(booking.BookingDetails)),
		}

		for j, detail := range booking.BookingDetails {
			resp.Data[i].Detail[j] = bookingdto.DetailBookingHistory{
				GuestName:        detail.Guest,
				AgentName:        booking.AgentName,
				HotelName:        detail.DetailRooms.HotelName,
				Additional:       detail.BookingDetailAdditionalName,
				SubBookingID:     detail.SubBookingID,
				BookingStatus:    detail.BookingStatus,
				PaymentStatus:    detail.PaymentStatus,
				CancellationDate: detail.DetailRooms.CancelledDate,
			}
		}
	}

	return resp, nil
}
