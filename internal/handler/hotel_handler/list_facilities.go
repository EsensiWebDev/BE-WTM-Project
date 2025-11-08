package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListFacilities godoc
// @Summary List Facilities
// @Description Retrieve a list of facilities with pagination.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter facilities by name"
// @Success 200 {object} response.ResponseWithPagination{data=[]string} "Successfully retrieved list of facilities"
// @Security BearerAuth
// @Router /hotels/facilities [get]
func (hh *HotelHandler) ListFacilities(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListFacilitiesRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	resp, err := hh.hotelUsecase.ListFacilities(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching facilities:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of facilities")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of facilities"

	var facilities []string
	if resp != nil {
		facilities = resp.Facilities
		if len(resp.Facilities) == 0 {
			message = "No facilities found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, facilities, message, pagination)
}
