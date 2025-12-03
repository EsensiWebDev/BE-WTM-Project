package booking_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// CancelBooking godoc
// @Summary      Cancel Booking
// @Description  Cancel an existing booking
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        sub_booking_id  path      string  true  "Booking ID"
// @Success      200         {object}  response.Response   "Successfully cancelled booking"
// @Security     BearerAuth
// @Router       /bookings/{sub_booking_id}/cancel [post]
func (bh *BookingHandler) CancelBooking(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.CancelBookingRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(ctx, "Failed to bind json parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := bh.bookingUsecase.CancelBooking(ctx, &req); err != nil {
		logger.Error(ctx, "Error cancelling booking:", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to cancel booking: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully cancelled booking")
}
