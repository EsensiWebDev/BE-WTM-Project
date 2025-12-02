package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// CreateHotel godoc
// @Summary Create Hotel
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Hotel name"
// @Param photos formData []file true "Hotel photos (multiple allowed)" collectionFormat(multi)
// @Param sub_district formData string true "Sub-district location"
// @Param district formData string true "District location"
// @Param province formData string true "Province location"
// @Param email formData string true "Contact email"
// @Param description formData string false "Hotel description"
// @Param rating formData int false "Hotel rating (1–5)"
// @Param nearby_places formData string false "Nearby places as JSON string. Example: /example_nearby_places "
// @Param facilities formData []string false "Facilities (multiple allowed)" collectionFormat(multi)
// @Param social_medias formData string false "Social media links as JSON string. Example: /example_social_medias "
// @Success 200 {object} response.ResponseWithData{data=hoteldto.CreateHotelResponse} "Successfully created hotel"
// @Router /hotels [post]
// @Security BearerAuth
func (hh *HotelHandler) CreateHotel(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.CreateHotelRequest
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

	resp, err := hh.hotelUsecase.CreateHotel(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Failed to create hotel", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to create hotel")
		return
	}

	response.Success(c, resp, "Successfully created hotel")
	return
}
