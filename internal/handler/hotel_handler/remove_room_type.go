package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemoveRoomType godoc
// @Summary Remove Room Type
// @Description Remove a room type by its Id.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param id path string true "Room Type Id"
// @Success 200 {object} response.Response "Successfully removed room type"
// @Security BearerAuth
// @Router /hotels/room-types/{id} [delete]
func (hh *HotelHandler) RemoveRoomType(c *gin.Context) {
	ctx := c.Request.Context()

	roomTypeID := c.Param("id")
	if roomTypeID == "" {
		logger.Error(ctx, "Hotel Id is required")
		response.Error(c, http.StatusBadRequest, "Hotel Id is required")
		return
	}

	// Convert hotelID to uint
	roomTypeIDint, err := utils.StringToUint(roomTypeID)
	if err != nil {
		logger.Error(ctx, "Invalid hotel Id format", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid hotel Id format")
		return
	}

	err = hh.hotelUsecase.RemoveRoomType(ctx, roomTypeIDint)
	if err != nil {
		logger.Error(ctx, "Error removing room type", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove room type")
		return
	}

	response.Success(c, nil, "Successfully removed room type")
}
