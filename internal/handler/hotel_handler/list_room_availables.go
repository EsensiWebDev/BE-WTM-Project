package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListRoomAvailable godoc
// @Summary List Room Available
// @Description Retrieve a list of available rooms based on the provided filters.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param hotel_id query uint true "Hotel Id to filter available rooms"
// @Param month query string true "Month to filter available rooms (format: YYYY-MM)"
// @Success 200 {object} response.ResponseWithData{data=[]hoteldto.RoomAvailable} "Successfully retrieved list of available rooms"
// @Security BearerAuth
// @Router /hotels/room-available [get]
func (hh *HotelHandler) ListRoomAvailable(c *gin.Context) {
	ctx := c.Request.Context()

	var req *hoteldto.ListRoomAvailableRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
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

	resp, err := hh.hotelUsecase.ListRoomAvailable(ctx, req)
	if err != nil {
		logger.Error(ctx, "Error fetching room available:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of room available")
		return
	}

	message := "Successfully retrieved list of available rooms"

	var roomAvailable []hoteldto.RoomAvailable
	if resp != nil {
		roomAvailable = resp.RoomAvailable
		if len(resp.RoomAvailable) == 0 {
			message = "No room available found"
		}
	}

	response.Success(c, roomAvailable, message)
}
