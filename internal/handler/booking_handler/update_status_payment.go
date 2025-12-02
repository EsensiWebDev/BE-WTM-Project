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

// UpdateStatusPayment godoc
// @Summary      Update payment status
// @Description  Update the status of a payment
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.UpdateStatusRequest true "Update status payment request"
// @Success      200 {object} response.Response "Successfully updated payment status"
// @Security     BearerAuth
// @Router       /bookings/payment-status [post]
func (bh *BookingHandler) UpdateStatusPayment(c *gin.Context) {
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

	if err := bh.bookingUsecase.UpdateStatusBooking(ctx, &req, constant.ConstPayment); err != nil {
		logger.Error(ctx, "Failed to update status payment", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update status payment")
		return
	}

	response.Success(c, nil, "Successfully updated payment status")

}
