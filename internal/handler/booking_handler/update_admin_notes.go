package booking_handler

import (
	"net/http"

	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UpdateAdminNotes godoc
// @Summary      Update admin notes for a booking detail
// @Description  Update the admin_notes field for a specific booking detail (sub-booking). This note will be visible to agents.
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        request body bookingdto.UpdateAdminNotesRequest true "Update admin notes request"
// @Success      200 {object} response.Response "Admin notes have been successfully updated"
// @Security     BearerAuth
// @Router       /bookings/admin-notes [post]
func (bh *BookingHandler) UpdateAdminNotes(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.UpdateAdminNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request for update admin notes:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Validation error for update admin notes", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := bh.bookingUsecase.UpdateAdminNotes(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update admin notes", err.Error())
		// Hide internal error details from client
		response.Error(c, http.StatusInternalServerError, "Failed to update admin notes")
		return
	}

	response.Success(c, nil, "Admin notes have been successfully updated")
}
