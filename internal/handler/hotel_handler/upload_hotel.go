package hotel_handler

import (
	"fmt"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// UploadHotel godoc
// @Summary Upload hotel
// @Description Upload hotel
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {object} response.Response "Successfully uploaded data hotel"
// @Router /hotels/upload [post]
// @Security BearerAuth
func (hh *HotelHandler) UploadHotel(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.UploadHotelRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind request payload", err)
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

	success, err := hh.hotelUsecase.UploadHotel(ctx, &req)
	if err != nil {
		if success {
			logger.Warn(ctx, "Partial upload success with errors", err.Error())
			response.Success(c, http.StatusOK, err.Error())
			return
		}
		logger.Error(ctx, "Failed to upload hotel", err.Error())
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to upload hotel: %s", err.Error()))
		return
	}

	response.Success(c, http.StatusOK, "Hotel uploaded successfully")
}
