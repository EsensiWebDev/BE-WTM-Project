package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListCart godoc
// @Summary      List cart items
// @Description  Retrieve all items in the user's cart
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Success 200 {object} response.ResponseWithPagination{data=bookingdto.ListCartResponse} "Successfully retrieved cart items"
// @Security BearerAuth
// @Router       /bookings/cart [get]
func (bh *BookingHandler) ListCart(c *gin.Context) {
	ctx := c.Request.Context()

	cart, err := bh.bookingUsecase.ListCart(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to list cart", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve cart")
		return
	}

	response.Success(c, cart, "Successfully retrieved cart items")
}
