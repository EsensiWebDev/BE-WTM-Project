package booking_handler

import "wtm-backend/internal/domain"

type BookingHandler struct {
	bookingUsecase domain.BookingUsecase
}

func NewBookingHandler(bookingUsecase domain.BookingUsecase) *BookingHandler {
	return &BookingHandler{
		bookingUsecase: bookingUsecase,
	}
}
