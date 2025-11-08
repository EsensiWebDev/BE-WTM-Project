package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListBookingHistory godoc
// @Summary      List Booking History
// @Description  List booking history with pagination and filtering options
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        page         query  int    false "Page number for pagination" default(1)
// @Param        limit        query  int    false "Number of items per page" default(10)
// @Param        search 	  query  string false "Search by booking code or guest name"
// @Param        search_by    query  string false "Search by guest name or booking code" Enums(guest_name, booking_id)
// @Param        status_booking_id query  int    false "Filter by booking status Id"
// @Param        status_payment_id query  int    false "Filter by payment status Id"
// @Success      200          {object} response.ResponseWithPagination{data=[]bookingdto.DataBookingHistory} "Successfully retrieved booking history"
// @Security     BearerAuth
// @Router       /bookings/history [get]
func (bh *BookingHandler) ListBookingHistory(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.ListBookingHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Failed to bind ListBookingHistoryRequest", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := bh.bookingUsecase.ListBookingHistory(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Failed to list booking history", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list booking history")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved booking history"

	var bookings []bookingdto.DataBookingHistory
	if resp != nil {
		bookings = resp.Data
		if len(bookings) == 0 {
			message = "No booking history found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, bookings, message, pagination)
}
