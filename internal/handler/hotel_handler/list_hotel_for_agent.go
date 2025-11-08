package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListHotelsForAgent godoc
// @Summary List Hotels for Agent
// @Description Retrieve a list of hotels for a specific agent with pagination and filtering options.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param limit query int false "Number of hotels to return per page (default: 10)"
// @Param page query int false "Page number to retrieve (default: 1)"
// @Param search query string false "Search term to filter hotels by name or description"
// @Param rating query array false "Filter hotels by rating (e.g., 1,2,3,4,5)"
// @Param district query array false "Filter hotels by district (e.g., 'Jakarta', 'Bandung')"
// @Param range_price_min query float64 false "Minimum price range for filtering hotels"
// @Param range_price_max query float64 false "Maximum price range for filtering hotels"
// @Param total_rooms query array false "Filter hotels by total number of rooms (e.g., 1,2,3,4,5)"
// @Param bed_type_id query array false "Filter hotels by bed type Id (e.g., 1,2,3)"
// @Success 200 {object} response.ResponseWithPagination{data=hoteldto.ListHotelForAgentResponse} "Successfully retrieved list of hotels for agent"
// @Security BearerAuth
// @Router /hotels/agent [get]
func (hh *HotelHandler) ListHotelsForAgent(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListHotelForAgentRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := hh.hotelUsecase.ListHotelsForAgent(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching hotels for agent:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of hotels for agent")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of hotels for agent"

	if resp != nil {
		if len(resp.Hotels) == 0 {
			message = "No users found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, resp, message, pagination)
}
