package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListStatusBooking godoc
// @Summary List Booking Statuses
// @Description Retrieve a list of booking statuses.
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseWithData{data=[]entity.StatusBooking} "Successfully retrieved list of booking statuses"
// @Security BearerAuth
// @Router /bookings/booking-status [get]
func (bh *BookingHandler) ListStatusBooking(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := bh.bookingUsecase.ListStatusBooking(ctx)
	if err != nil {
		logger.Error(ctx, "Error listing booking statuses:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list booking statuses")
		return
	}

	if resp == nil || len(resp.Data) == 0 {
		logger.Error(ctx, "No booking statuses found")
		response.Success(c, http.StatusInternalServerError, "No booking statuses found")
		return
	}

	response.Success(c, resp.Data, "Successfully retrieved list of booking statuses")
	return
}
