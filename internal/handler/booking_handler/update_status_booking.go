package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateStatusBooking godoc
// @Summary      Update booking status
// @Description  Update the status of a booking
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.UpdateStatusRequest true "Update status booking request"
// @Success      200 {object} response.Response "Successfully updated booking status"
// @Security     BearerAuth
// @Router       /bookings/booking-status [post]
func (bh *BookingHandler) UpdateStatusBooking(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := bh.bookingUsecase.UpdateStatusBooking(ctx, &req, constant.ConstBooking); err != nil {
		logger.Error(ctx, "Failed to update status booking", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, nil, "Successfully updated booking status")

}
