package booking_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateCartAdditionalNotes godoc
// @Summary      Update additional notes for a cart detail item
// @Description  Update the admin/agent-only additional_notes field for a specific sub-cart (booking_detail) item
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.UpdateCartAdditionalNotesRequest true "Update additional notes request"
// @Success 200 {object} response.Response "Additional notes have been successfully updated"
// @Security BearerAuth
// @Router       /bookings/cart/sub-notes [post]
func (bh *BookingHandler) UpdateCartAdditionalNotes(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.UpdateCartAdditionalNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request for update cart additional notes:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error for update cart additional notes", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := bh.bookingUsecase.UpdateCartAdditionalNotes(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update cart additional notes", err.Error())
		// Hide internal error details from client
		response.Error(c, http.StatusInternalServerError, "Failed to update additional notes")
		return
	}

	response.Success(c, nil, "Additional notes have been successfully updated")
}



