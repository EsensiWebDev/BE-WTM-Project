package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateRoomAvailable godoc
// @Summary Update Room Available
// @Tags Hotel
// @Accept json
// @Produce json
// @Param request body hoteldto.UpdateRoomAvailableRequest true "Update Room Available Request"
// @Success 200 {object} response.Response "Successfully updated room available"
// @Router /hotels/room-available [put]
// @Security BearerAuth
func (hh *HotelHandler) UpdateRoomAvailable(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.UpdateRoomAvailableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind request", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to bind request")
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

	if err := hh.hotelUsecase.UpdateRoomAvailable(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update room available", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update room available")
		return
	}

	response.Success(c, nil, "Successfully updated room available")

}
