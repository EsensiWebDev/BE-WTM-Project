package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListProvinces godoc
// @Summary List Provinces
// @Description Retrieve a list of provinces where hotels are located.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} response.ResponseWithPagination{data=[]string} "Successfully retrieved list of provinces"
// @Router /hotels/provinces [get]
// @Security BearerAuth
func (hh *HotelHandler) ListProvinces(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListProvincesRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Failed to bind ListProvincesRequest", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := hh.hotelUsecase.ListProvinces(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Failed to list provinces", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to list provinces")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of provinces"

	var provinces []string
	if resp != nil {
		provinces = resp.Provinces
		if len(provinces) == 0 {
			message = "No provinces found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, provinces, message, pagination)
}
