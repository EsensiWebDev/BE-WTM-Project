package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemoveHotel godoc
// @Summary Remove Hotel
// @Description Remove a hotel by its Id.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param id path string true "Hotel Id"
// @Success 200 {object} response.Response "Successfully removed hotel"
// @Security BearerAuth
// @Router /hotels/{id} [delete]
func (hh *HotelHandler) RemoveHotel(c *gin.Context) {
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

	err = hh.hotelUsecase.RemoveHotel(ctx, hotelIDUint)
	if err != nil {
		logger.Error(ctx, "Error removing hotel", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove hotel")
		return
	}

	response.Success(c, nil, "Successfully removed hotel")
}
