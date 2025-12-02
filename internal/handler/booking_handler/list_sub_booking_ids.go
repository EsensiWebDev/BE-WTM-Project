package booking_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bookingdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListSubBookingIDs godoc
// @Summary      List Sub Booking IDs
// @Description  List sub booking IDs for a given booking ID associated with the authenticated agent
// @Tags         Booking
// @Accept       json
// @Produce      json
// @Param        booking_id   path      string  true  "Booking ID"
// @Success      200          {object}  response.Response{data=[]string} "Successfully retrieved sub booking IDs"
// @Security     BearerAuth
// @Router       /bookings/{booking_id}/sub-ids [get]
func (bh *BookingHandler) ListSubBookingIDs(c *gin.Context) {
	ctx := c.Request.Context()

	var req bookingdto.ListSubBookingIDsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(ctx, "Failed to bind path parameter ", err.Error())
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

	resp, err := bh.bookingUsecase.ListSubBookingIDs(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching sub booking IDs:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch sub booking IDs")
		return
	}

	message := "Successfully retrieved sub booking IDs"

	var subBookingIDs []string
	if resp != nil {
		subBookingIDs = resp.SubBookingIDs
		if len(subBookingIDs) == 0 {
			message = "No sub booking IDs found"
		}
	}

	response.Success(c, subBookingIDs, message)
}
