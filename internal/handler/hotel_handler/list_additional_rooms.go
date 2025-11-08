package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListAdditionalRooms godoc
// @Summary List Additional Rooms
// @Description Retrieve a list of additional rooms with pagination.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter additional rooms by name"
// @Success 200 {object} response.ResponseWithPagination{data=[]string} "Successfully retrieved list of additional rooms"
// @Security BearerAuth
// @Router /hotels/additional-rooms [get]
func (hh *HotelHandler) ListAdditionalRooms(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListAdditionalRoomsRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := hh.hotelUsecase.ListAdditionalRooms(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching additional rooms:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of additional rooms")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of additional rooms"

	var additionalRooms []string
	if resp != nil {
		additionalRooms = resp.AdditionalRooms
		if len(resp.AdditionalRooms) == 0 {
			message = "No facilities found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, additionalRooms, message, pagination)
}
