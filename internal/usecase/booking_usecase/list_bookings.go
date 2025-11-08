package booking_usecase

import (
	"context"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookings(ctx context.Context, req *bookingdto.ListBookingsRequest) (resp *bookingdto.ListBookingsResponse, err error) {
	filterReq := filter.BookingFilter{}
	filterReq.PaginationRequest = req.PaginationRequest
	filterReq.StatusBookingID = req.StatusBookingID

	bookings, total, err := bu.bookingRepo.GetBookings(ctx, &filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get bookings", err.Error())
		return nil, err
	}

	resp = &bookingdto.ListBookingsResponse{
		Total: total,
		Data:  make([]bookingdto.DataBooking, len(bookings)),
	}

	for i, booking := range bookings {
		resp.Data[i] = bookingdto.DataBooking{
			BookingID:    booking.ID,
			GuestName:    booking.Guests,
			AgentName:    booking.AgentName,
			AgentCompany: booking.AgentCompanyName,
			//GroupPromo:    booking.GroupPromo,
			BookingCode:   booking.BookingCode,
			BookingStatus: booking.BookingStatus,
			PaymentStatus: booking.PaymentStatus,
			Detail:        make([]bookingdto.DetailBooking, len(booking.BookingDetails)),
		}

		for j, detail := range booking.BookingDetails {
			resp.Data[i].Detail[j] = bookingdto.DetailBooking{
				GuestName:     detail.Guest,
				HotelName:     detail.DetailRooms.HotelName,
				Additional:    detail.BookingDetailAdditionalName,
				SubBookingID:  detail.SubBookingID,
				BookingStatus: booking.BookingStatus,
				PaymentStatus: booking.PaymentStatus,
				IsAPI:         detail.DetailRooms.IsAPI,
				PromoID:       detail.PromoID,
				PromoCode:     detail.DetailPromos.PromoCode,
			}
			if detail.StatusBookingID != constant.StatusBookingRejectedID {
				resp.Data[i].Detail[j].CancelledDate = detail.UpdatedAt.Format("2006-01-02")
			}
		}
	}

	return resp, nil
}
