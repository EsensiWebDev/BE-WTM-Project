package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateRoomType godoc
// @Summary Update Room Type
// @Description Update the details of a room type by Id.
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Room Type Id"
// @Param name formData string true "Room Type Name"
// @Param photos formData []file false "Room Type Photos (multiple allowed)" collectionFormat(multi)
// @Param without_breakfast formData string true "Price without breakfast as JSON string. Example: /example_room_options?room_options=without_breakfast "
// @Param with_breakfast formData string true "Price with breakfast as JSON string. Example: /example_room_options?room_options=with_breakfast "
// @Param room_size formData number true "Room Size in square meters"
// @Param max_occupancy formData int true "Maximum Occupancy"
// @Param bed_types formData []string true "Bed Types (multiple allowed)" collectionFormat(multi)
// @Param is_smoking_room formData bool true "Is Smoking Room"
// @Param additional formData string false "Additional amenities as JSON string. Example: /example_additional_features"
// @Param description formData string false "Room Type Description"
// @Param unchanged_room_photos formData []string false "Unchanged room photos (multiple allowed)" collectionFormat(multi)
// @Param unchanged_additions_ids formData []int false "Unchanged addition IDs (multiple allowed)" collectionFormat(multi)
// @Success 200 {object} response.Response "Successfully updated room type"
// @Security BearerAuth
// @Router /hotels/room-types/{id} [put]
func (hh *HotelHandler) UpdateRoomType(c *gin.Context) {
	ctx := c.Request.Context()

	roomTypeID := c.Param("id")
	if roomTypeID == "" {
		logger.Error(ctx, "Room Type Id is required")
		response.Error(c, http.StatusBadRequest, "Room Type Id is required")
		return
	}

	roomTypeIDUint, err := utils.StringToUint(roomTypeID)
	if err != nil {
		logger.Error(ctx, "Failed to convert Room Type Id to uint", err.Error())
		response.Error(c, http.StatusBadRequest, "Failed to convert Room Type Id to uint")
		return
	}

	var req hoteldto.UpdateRoomTypeRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind UpdateRoomTypeRequest", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.RoomTypeID = roomTypeIDUint

	if err := hh.hotelUsecase.UpdateRoomType(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update room type", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update room type")
		return
	}

	response.Success(c, nil, "Successfully updated room type")

}
