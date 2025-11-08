package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// CheckOutCart godoc
// @Summary      Checkout cart
// @Description  Finalize the cart by submitting guest data and changing booking status to "in review"
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.CheckOutCartRequest true "Checkout cart request"
// @Success      200 {object} response.Response "Successfully checked out cart"
// @Failure      400 {object} response.Response "Invalid request payload"
// @Failure      500 {object} response.Response "Internal server error"
// @Security     BearerAuth
// @Router       /bookings/checkout [post]
func (bh *BookingHandler) CheckOutCart(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.CheckOutCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind JSON body", err.Error())
		response.Error(c, http.StatusBadRequest, "Failed to bind JSON body")
		return
	}

	if err := bh.bookingUsecase.CheckOutCart(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to check out cart", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to check out cart")
		return
	}

	response.Success(c, nil, "Successfully checked out cart")
}
