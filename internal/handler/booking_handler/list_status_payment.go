package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListStatusPayment godoc
// @Summary List Payment Statuses
// @Description Retrieve a list of payment statuses.
// @Tags Booking
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseWithData{data=[]entity.StatusPayment} "Successfully retrieved list of payment statuses"
// @Security BearerAuth
// @Router /bookings/payment-status [get]
func (bh *BookingHandler) ListStatusPayment(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := bh.bookingUsecase.ListStatusPayment(ctx)
	if err != nil {
		logger.Error(ctx, "Error listing payment statuses:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list payment statuses")
		return
	}

	if resp == nil || len(resp.Data) == 0 {
		logger.Error(ctx, "No payment statuses found")
		response.Success(c, http.StatusInternalServerError, "No payment statuses found")
		return
	}

	response.Success(c, resp.Data, "Successfully retrieved list of payment statuses")
	return
}
