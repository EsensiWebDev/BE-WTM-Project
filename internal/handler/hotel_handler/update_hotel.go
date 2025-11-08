package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateHotel godoc
// @Summary Update Hotel
// @Description Update hotel information by Id
// @Tags Hotel
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Hotel ID"
// @Param name formData string true "Hotel name"
// @Param photos formData []file false "Hotel photos (multiple allowed)" collectionFormat(multi)
// @Param sub_district formData string true "Sub-district location"
// @Param district formData string true "District location"
// @Param province formData string true "Province location"
// @Param email formData string true "Contact email"
// @Param description formData string false "Hotel description"
// @Param rating formData int false "Hotel rating (1–5)"
// @Param nearby_places formData string false "Nearby places as JSON string. Example: /example_nearby_places "
// @Param facilities formData []string false "Facilities (multiple allowed)" collectionFormat(multi)
// @Param social_medias formData string false "Social media links as JSON string. Example: /example_social_medias "
// @Param unchanged_hotel_photos formData []string false "Unchanged hotel photos (multiple allowed)" collectionFormat(multi)
// @Param unchanged_nearby_place_ids formData []int false "Unchanged nearby place IDs (multiple allowed)" collectionFormat(multi)
// @Success 200 {object} response.Response "Successfully updated hotel"
// @Router /hotels/{id} [put]
// @Security BearerAuth
func (hh *HotelHandler) UpdateHotel(c *gin.Context) {
	ctx := c.Request.Context()

	hotelID := c.Param("id")
	if hotelID == "" {
		logger.Error(ctx, "Hotel Id is required")
		response.Error(c, http.StatusBadRequest, "Hotel Id is required")
		return
	}

	// Convert hotelID to uint
	hotelIDUint, err := utils.StringToUint(hotelID)
	if err != nil {
		logger.Error(ctx, "Invalid hotel Id format", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid hotel Id format")
		return
	}

	var req hoteldto.UpdateHotelRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind JSON body", err.Error())
		response.Error(c, http.StatusBadRequest, "Failed to bind JSON body")
		return
	}

	req.HotelID = hotelIDUint

	if err := hh.hotelUsecase.UpdateHotel(ctx, &req); err != nil {
		logger.Error(ctx, "Failed to update hotel", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to update hotel")
		return
	}

	response.Success(c, nil, "Successfully updated hotel")
	return
}
