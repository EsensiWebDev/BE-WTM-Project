package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AddToCart godoc
// @Summary      Add item to cart
// @Description  Add an item to the user's cart
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.AddToCartRequest true "Add to cart request"
// @Success 200 {object} response.Response "Successfully added to cart"
// @Security BearerAuth
// @Router       /bookings/cart [post]
func (bh *BookingHandler) AddToCart(c *gin.Context) {

	ctx := c.Request.Context()

	var req bookingdto.AddToCartRequest

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

	if err := bh.bookingUsecase.AddToCart(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to add to cart", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to add to cart")
		return
	}

	response.Success(c, nil, "Successfully added to cart")
}
