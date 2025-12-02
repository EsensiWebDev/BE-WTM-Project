package booking_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// CheckOutCart godoc
// @Summary      Checkout cart
// @Description  Finalize the cart by submitting guest data and changing booking status to "in review"
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response{data=[]bookingdto.DataInvoice} "Successfully checked out cart"
// @Security     BearerAuth
// @Router       /bookings/checkout [post]
func (bh *BookingHandler) CheckOutCart(c *gin.Context) {
	ctx := c.Request.Context()

	dataCheckout, err := bh.bookingUsecase.CheckOutCart(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to check out cart", err.Error())
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to check out cart: %s", err.Error()))
		return
	}
	if len(dataCheckout.Invoice) == 0 {
		logger.Error(ctx, "Invoice data is empty")
		response.Error(c, http.StatusBadRequest, "Invoice data is empty")
		return
	}

	response.Success(c, dataCheckout.Invoice, "Successfully checked out cart")
}
