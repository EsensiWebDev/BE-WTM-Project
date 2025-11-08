package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// DetailHotel godoc
// @Summary Get Hotel Details by Id
// @Description Retrieve hotel details by Id.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param id path string true "Hotel Id"
// @Success 200 {object} response.ResponseWithData{data=hoteldto.DetailHotelResponse} "Successfully retrieved hotel details"
// @Security BearerAuth
// @Router /hotels/{id} [get]
func (hh *HotelHandler) DetailHotel(c *gin.Context) {
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

	hotel, err := hh.hotelUsecase.DetailHotel(ctx, hotelIDUint)
	if err != nil {
		logger.Error(ctx, "Error getting hotel by Id", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get hotel details")
		return
	}

	if hotel == nil {
		response.Error(c, http.StatusNotFound, "Hotel not found")
		return
	}

	response.Success(c, hotel, "Successfully retrieved hotel details")
}
