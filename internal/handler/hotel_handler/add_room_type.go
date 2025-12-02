package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AddRoomType godoc
// @Summary Add Room Type
// @Description Add a new room type to a hotel
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param hotel_id formData int true "Hotel ID"
// @Param name formData string true "Room type name"
// @Param photos formData []file true "Room type photos (multiple allowed)" collectionFormat(multi)
// @Param without_breakfast formData string false "Without breakfast details as JSON string. Example : /example_room_options?room_options=without_breakfast"
// @Param with_breakfast formData string false "With breakfast details as JSON string. Example : /example_room_options?room_options=with_breakfast"
// @Param room_size formData number false "Room size in square meters"
// @Param max_occupancy formData int true "Maximum occupancy"
// @Param bed_types formData []string true "Bed types (multiple allowed)" collectionFormat(multi)
// @Param is_smoking_room formData bool false "Is smoking room"
// @Param additional formData string false "Additional room features as JSON string. Example: /example_additional_features"
// @Param description formData string false "Room type description"
// @Success 200 {object} response.Response "Successfully added room type"
// @Router /hotels/room-types [post]
// @Security BearerAuth
func (hh *HotelHandler) AddRoomType(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.AddRoomTypeRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind AddRoomTypeRequest", err.Error())
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

	if err := hh.hotelUsecase.AddRoomType(ctx, req.HotelID, &req); err != nil {
		logger.Error(ctx, "Failed to add room type", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to add room type")
		return
	}

	response.Success(c, nil, "Successfully added room type")
	return
}
