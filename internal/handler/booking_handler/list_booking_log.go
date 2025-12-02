package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListBookingLog godoc
// @Summary      List Booking Log
// @Description  List booking logs with pagination and filtering options
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        page         query  int    false "Page number for pagination" default(1)
// @Param        limit        query  int    false "Number of items per page" default(10)
// @Param        search       query  string false "Search by booking code or guest name"
// @Param        booking_status_id query  int    false "Filter by booking status ID"
// @Param        payment_status_id query  int    false "Filter by payment status ID"
// @Param        confirm_date_from query  string false "Filter by confirm date from (YYYY-MM-DD)"
// @Param        confirm_date_to query  string false "Filter by confirm date to (YYYY-MM-DD)"
// @Param        check_in_date_from query  string false "Filter by check-in date from (YYYY-MM-DD)"
// @Param        check_in_date_to query  string false "Filter by check-in date to (YYYY-MM-DD)"
// @Param        check_out_date_from query  string false "Filter by check-out date from (YYYY-MM-DD)"
// @Param        check_out_date_to query  string false "Filter by check-out date to (YYYY-MM-DD)"
// @Success      200          {object} response.ResponseWithPagination{data=[]bookingdto.BookingLog} "Successfully retrieved booking logs"
// @Security     BearerAuth
// @Router       /bookings/logs [get]
func (bh *BookingHandler) ListBookingLog(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.ListBookingLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Failed to bind ListBookingLogRequest", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := bh.bookingUsecase.ListBookingLog(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Failed to list booking log", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list booking log")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved booking logs"

	var logs []bookingdto.BookingLog
	if resp != nil {
		logs = resp.Data
		if len(logs) == 0 {
			message = "No booking logs found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, logs, message, pagination)
}
