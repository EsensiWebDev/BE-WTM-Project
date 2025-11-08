package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListRoomTypes godoc
// @Summary List Room Types
// @Description Retrieve a list of room types for a specific hotel.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param hotel_id query uint true "Hotel Id to filter room types"
// @Success 200 {object} response.ResponseWithData{data=[]hoteldto.ListRoomType} "Successfully retrieved list of room types"
// @Security BearerAuth
// @Router /hotels/room-types [get]
func (hh *HotelHandler) ListRoomTypes(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListRoomTypeRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
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

	resp, err := hh.hotelUsecase.ListRoomTypes(ctx, req.HotelID)
	if err != nil {
		logger.Error(ctx, "Error getting data list room type", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list room type")
	}

	message := "Successfully retrieved list of room types"

	var roomTypes []hoteldto.ListRoomType
	if resp != nil {
		roomTypes = resp.RoomTypes
		if len(resp.RoomTypes) == 0 {
			message = "No room types found"
		}
	}

	response.Success(c, roomTypes, message)

}
