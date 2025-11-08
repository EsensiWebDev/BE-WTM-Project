package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemoveFromCart godoc
// @Summary      Remove item from cart
// @Description  Remove an item from the user's cart
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        id path string true "Booking Detail Id"
// @Success 200 {object} response.Response "Successfully removed item from cart"
// @Security BearerAuth
// @Router       /bookings/cart/{id} [delete]
func (bh *BookingHandler) RemoveFromCart(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Booking Id is required")
		return
	}

	bookingDetailID, err := utils.StringToUint(id)
	if err != nil {
		logger.Error(ctx, "Invalid booking detail Id", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid booking detail Id")
		return
	}

	if err := bh.bookingUsecase.RemoveFromCart(ctx, bookingDetailID); err != nil {
		logger.Error(ctx, "Failed to remove from cart", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove from cart")
		return
	}

	response.Success(c, nil, "Successfully removed item from cart")
}
