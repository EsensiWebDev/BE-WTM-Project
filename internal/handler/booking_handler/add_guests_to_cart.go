package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AddGuestsToCart godoc
// @Summary      Add Guests to Cart
// @Description  Add guests to a booking cart
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request  body      bookingdto.AddGuestsToCartRequest  true  "Add Guests to Cart Request"
// @Success      200      {object}  response.Response   "Successfully added guests to cart"
// @Security     BearerAuth
// @Router       /bookings/cart/guests [post]
func (bh *BookingHandler) AddGuestsToCart(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.AddGuestsToCartRequest
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

	if err := bh.bookingUsecase.AddGuestsToCart(ctx, &req); err != nil {
		logger.Error(ctx, "Error adding guests to cart:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to add guests to cart")
		return
	}

	response.Success(c, nil, "Successfully added guests to cart")
}
