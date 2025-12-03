package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

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

	if err := hh.hotelUsecase.UploadHotel(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to upload hotel", err)
		response.Error(c, http.StatusInternalServerError, "Failed to upload hotel")
		return
	}

	response.Success(c, http.StatusOK, "Hotel uploaded successfully")
}
