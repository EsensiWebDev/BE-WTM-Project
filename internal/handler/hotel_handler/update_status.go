package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateStatus godoc
// @Summary Update Hotel Status
// @Description Update the status of a hotel by Id.
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param hotel_id formData int true "Hotel Id"
// @Param status formData bool true "Hotel status (true for approved, false for rejected)"
// @Success 200 {object} response.Response "Successfully updated hotel status"
// @Security BearerAuth
// @Router /hotels/status [put]
func (hh *HotelHandler) UpdateStatus(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.UpdateStatusRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind CreateHotelRequest", err.Error())
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

	if err := hh.hotelUsecase.UpdateStatus(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update hotel status", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update hotel status")
		return
	}

	response.Success(c, nil, "Successfully updated hotel status")
}
