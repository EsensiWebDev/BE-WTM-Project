package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListBookings godoc
// @Summary      List Bookings
// @Description  List bookings with pagination and filtering options
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        page         query  int    false "Page number for pagination" default(1)
// @Param        limit        query  int    false "Number of items per page" default(10)
// @Param        search       query  string false "Search by booking code or guest name"
// @Param		 booking_status_id query  int    false "Filter by booking status Id"
// @Param        payment_status_id query  int    false "Filter by payment status Id"
// @Success      200          {object} response.ResponseWithPagination{data=[]bookingdto.DataBooking} "Successfully retrieved bookings"
// @Security     BearerAuth
// @Router       /bookings [get]
func (bh *BookingHandler) ListBookings(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.ListBookingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Failed to bind ListBookingsRequest", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}
	resp, err := bh.bookingUsecase.ListBookings(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching bookings:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved bookings"

	var bookings []bookingdto.DataBooking
	if resp != nil {
		bookings = resp.Data
		if len(bookings) == 0 {
			message = "No bookings found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, bookings, message, pagination)
}
