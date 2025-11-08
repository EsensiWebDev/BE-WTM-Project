package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListAllBedTypes godoc
// @Summary List All Bed Types
// @Description Retrieve a list of all bed types with pagination.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword to filter bed types by name"
// @Success 200 {object} response.ResponseWithPagination{data=[]string} "Successfully retrieved list of all bed types"
// @Security BearerAuth
// @Router /hotels/bed-types [get]
func (hh *HotelHandler) ListAllBedTypes(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListAllBedTypesRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
	}

	resp, err := hh.hotelUsecase.ListAllBedTypes(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching all bed types:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of all bed types")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of all bed types"

	var allBedTypes []string
	if resp != nil {
		allBedTypes = resp.BedTypes
		if len(resp.BedTypes) == 0 {
			message = "No facilities found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, allBedTypes, message, pagination)
}
