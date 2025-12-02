package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AddGuestToSubCart godoc
// @Summary      Add/Change Guest to Sub Cart
// @Description  Add or change a guest to a specific sub cart within a booking
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request  body      bookingdto.AddGuestToSubCartRequest  true  "Add / Change Guest to Sub Cart Request"
// @Success      200      {object}  response.Response   "Successfully added / changed guest to sub cart"
// @Security     BearerAuth
// @Router       /bookings/cart/sub-guest [post]
func (bh *BookingHandler) AddGuestToSubCart(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.AddGuestToSubCartRequest
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

	if err := bh.bookingUsecase.AddGuestToSubCart(ctx, &req); err != nil {
		logger.Error(ctx, "Error adding / changing guest to sub cart:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to add / change guest to sub cart")
		return
	}

	response.Success(c, nil, "Successfully added / changed guest to sub cart")
}
