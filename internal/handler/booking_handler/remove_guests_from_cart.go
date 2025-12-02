package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemoveGuestsFromCart godoc
// @Summary      Remove Guests from Cart
// @Description  Remove guests from a booking cart
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request  body      bookingdto.RemoveGuestsFromCartRequest  true  "Remove Guests from Cart Request"
// @Success      200      {object}  response.Response   "Successfully removed guests from cart"
// @Security     BearerAuth
// @Router       /bookings/cart/guests [delete]
func (bh *BookingHandler) RemoveGuestsFromCart(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.RemoveGuestsFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	if err := bh.bookingUsecase.RemoveGuestsFromCart(ctx, &req); err != nil {
		logger.Error(ctx, "Error removing guests from cart:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove guests from cart")
		return
	}

	response.Success(c, nil, "Successfully removed guests from cart")
}
