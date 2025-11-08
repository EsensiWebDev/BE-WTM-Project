package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListStatusHotel godoc
// @Summary List Hotel Statuses
// @Description Retrieve a list of hotel statuses.
// @Tags Hotel
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseWithData{data=[]entity.StatusHotel} "Successfully retrieved list of hotel statuses"
// @Security BearerAuth
// @Router /hotels/statuses [get]
func (hh *HotelHandler) ListStatusHotel(c *gin.Context) {
	ctx := c.Request.Context()

	resp, err := hh.hotelUsecase.ListStatusHotel(ctx)
	if err != nil {
		logger.Error(ctx, "Error listing hotel statuses:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list hotel statuses")
		return
	}

	if resp == nil || len(resp.StatusHotel) == 0 {
		logger.Error(ctx, "No hotel statuses found")
		response.Success(c, http.StatusInternalServerError, "No hotel statuses found")
		return
	}

	response.Success(c, resp.StatusHotel, "Successfully retrieved list of hotel statuses")
	return
}
