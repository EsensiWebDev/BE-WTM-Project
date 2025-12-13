package booking_usecase

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

func (bu *BookingUsecase) ListBookings(ctx context.Context, req *bookingdto.ListBookingsRequest) (resp *bookingdto.ListBookingsResponse, err error) {
	filterReq := filter.BookingFilter{}
	filterReq.PaginationRequest = req.PaginationRequest
	filterReq.PaymentStatusID = req.PaymentStatusID
	filterReq.BookingStatusID = req.BookingStatusID

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
			BookingID:     booking.BookingCode,
			AgentName:     booking.AgentName,
			AgentCompany:  booking.AgentCompanyName,
			GroupPromo:    booking.PromoGroupAgent,
			PaymentStatus: booking.PaymentStatus,
			BookingStatus: booking.BookingStatus,
			Detail:        make([]bookingdto.DetailBooking, len(booking.BookingDetails)),
		}

		for j, detail := range booking.BookingDetails {
			resp.Data[i].Detail[j] = bookingdto.DetailBooking{
				HotelName:     detail.DetailRooms.HotelName,
				Additional:    detail.BookingDetailAdditionalName,
				SubBookingID:  detail.SubBookingID,
				BookingStatus: detail.BookingStatus,
				PaymentStatus: detail.PaymentStatus,
				IsAPI:         detail.DetailRooms.IsAPI,
				PromoCode:     detail.DetailPromos.PromoCode,
			}
			var receiptUrl string
			if detail.ReceiptUrl != "" {
				bucketName := fmt.Sprintf("%s-%s", constant.ConstBooking, constant.ConstPrivate)
				receiptUrl, err = bu.fileStorage.GetFile(ctx, bucketName, detail.ReceiptUrl)
				if err != nil {
					logger.Error(ctx, "ListHotelsForAgent", err.Error())
				}
				resp.Data[i].Detail[j].Receipt = receiptUrl
				resp.Data[i].Receipts = append(resp.Data[i].Receipts, receiptUrl)
			}
			if detail.Invoice != nil {
				resp.Data[i].Detail[j].Invoice = bookingdto.DataInvoice{
					InvoiceNumber: detail.Invoice.InvoiceCode,
					DetailInvoice: detail.Invoice.DetailInvoice,
					InvoiceDate:   detail.Invoice.CreatedAt.Format(time.DateOnly),
					Receipt:       receiptUrl,
				}
			}
			if detail.StatusBookingID != constant.StatusBookingRejectedID {
				resp.Data[i].Detail[j].CancelledDate = detail.ApprovedAt.Format("2006-01-02")
			}
			if strings.TrimSpace(detail.Guest) != "" {
				resp.Data[i].Detail[j].GuestName = detail.Guest
				resp.Data[i].GuestName = append(resp.Data[i].GuestName, detail.Guest)
			}
		}
	}

	return resp, nil
}
