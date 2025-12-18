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
		PaginationRequest: req.PaginationRequest,
		AgentID:           agentID,
		BookingStatusID:   req.StatusBookingID,
		PaymentStatusID:   req.StatusPaymentID,
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
		var statusBooking []string
		var statusPayment []string
		resp.Data[i] = bookingdto.DataBookingHistory{
			BookingID:     booking.ID,
			BookingCode:   booking.BookingCode,
			BookingStatus: booking.BookingStatus,
			PaymentStatus: booking.PaymentStatus,
			Detail:        make([]bookingdto.DetailBookingHistory, len(booking.BookingDetails)),
		}

		for j, detail := range booking.BookingDetails {
			// Parse other preferences from comma-separated snapshot (same logic as cart)
			var otherPrefs []string
			if strings.TrimSpace(detail.OtherPreferences) != "" {
				for _, p := range strings.Split(detail.OtherPreferences, ",") {
					if name := strings.TrimSpace(p); name != "" {
						otherPrefs = append(otherPrefs, name)
					}
				}
			}

			// Map detailed additional services (with price, category, pax, etc.)
			var additionalServices []bookingdto.BookingHistoryAdditional
			var totalAdditionalPrice float64
			if len(detail.BookingDetailsAdditional) > 0 {
				additionalServices = make([]bookingdto.BookingHistoryAdditional, 0, len(detail.BookingDetailsAdditional))
				for _, add := range detail.BookingDetailsAdditional {
					additionalService := bookingdto.BookingHistoryAdditional{
						Name:       add.NameAdditional,
						Category:   add.Category,
						Price:      add.Price,
						Pax:        add.Pax,
						IsRequired: add.IsRequired,
					}
					additionalServices = append(additionalServices, additionalService)
					// Calculate total additional price
					if add.Category == constant.AdditionalServiceCategoryPrice && add.Price != nil {
						totalAdditionalPrice += *add.Price
					}
				}
			}

			// Calculate nights
			nights := int(detail.CheckOutDate.Sub(detail.CheckInDate).Hours() / 24)
			if nights < 1 {
				nights = 1
			}

			// Get currency from booking detail (snapshot at booking time)
			bookingCurrency := detail.Currency
			if bookingCurrency == "" {
				bookingCurrency = "IDR" // Default fallback
			}

			// Room price per night (already includes promo if applied)
			roomPricePerNight := detail.Price / float64(nights)
			if nights == 0 {
				roomPricePerNight = detail.Price
			}

			// Total price = room price + additional services
			totalPrice := detail.Price + totalAdditionalPrice

			resp.Data[i].Detail[j] = bookingdto.DetailBookingHistory{
				GuestName:          detail.Guest,
				AgentName:          booking.AgentName,
				HotelName:          detail.DetailRooms.HotelName,
				RoomTypeName:       detail.DetailRooms.RoomTypeName,
				IsBreakfast:        detail.RoomPrice.IsBreakfast,
				BedType:            detail.BedType,
				RoomPrice:          roomPricePerNight,
				TotalPrice:         totalPrice,
				Currency:           bookingCurrency,
				CheckInDate:        detail.CheckInDate.Format(time.DateOnly),
				CheckOutDate:       detail.CheckOutDate.Format(time.DateOnly),
				Additional:         detail.BookingDetailAdditionalName, // Keep for backward compatibility
				AdditionalServices: additionalServices,                 // Detailed additional services
				OtherPreferences:   otherPrefs,
				SubBookingID:       detail.SubBookingID,
				BookingStatus:      detail.BookingStatus,
				PaymentStatus:      detail.PaymentStatus,
				CancellationDate:   detail.DetailRooms.CancelledDate,
				AdditionalNotes:    detail.AdditionalNotes,
				AdminNotes:         detail.AdminNotes,
			}
			statusBooking = append(statusBooking, detail.BookingStatus)
			statusPayment = append(statusPayment, detail.PaymentStatus)
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
				invoice := bookingdto.DataInvoice{
					InvoiceNumber: detail.Invoice.InvoiceCode,
					DetailInvoice: detail.Invoice.DetailInvoice,
					InvoiceDate:   detail.Invoice.CreatedAt.Format(time.DateOnly),
					Receipt:       receiptUrl,
				}
				resp.Data[i].Detail[j].Invoice = invoice
				resp.Data[i].Invoices = append(resp.Data[i].Invoices, invoice)
			}
			if detail.Guest != "" {
				resp.Data[i].GuestName = append(resp.Data[i].GuestName, detail.Guest)
			}
		}
		resp.Data[i].BookingStatus = bu.summaryStatus(statusBooking, constant.ConstBooking)
		resp.Data[i].PaymentStatus = bu.summaryStatus(statusPayment, constant.ConstPayment)

	}

	return resp, nil
}
