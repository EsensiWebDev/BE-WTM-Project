package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListBookingIDs godoc
// @Summary      List Booking IDs
// @Description  List booking IDs with pagination
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        page         query  int    false "Page number for pagination" default(1)
// @Param        limit        query  int    false "Number of items per page" default(10)
// @Success      200          {object} response.ResponseWithPagination{data=[]string} "Successfully retrieved booking IDs"
// @Security     BearerAuth
// @Router       /bookings/ids [get]
func (bh *BookingHandler) ListBookingIDs(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.ListBookingIDsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Error(ctx, "Failed to bind query parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := bh.bookingUsecase.ListBookingIDs(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching booking IDs:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch booking IDs")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved booking IDs"

	var bookingIDs []string
	if resp != nil {
		bookingIDs = resp.BookingIDs
		if len(bookingIDs) == 0 {
			message = "No booking IDs found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, bookingIDs, message, pagination)
}
