package hotel_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/hoteldto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// ListHotels godoc
// @Summary List Hotels
// @Description Retrieve a paginated list of hotels using query parameters.
// @Tags Hotel
// @Accept json
// @Produce json
// @Param is_api query bool false "Filter hotels by API status"
// @Param region query string false "Filter hotels by region"
// @Param status_id query int false "Filter hotels by status_id"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search keyword to filter hotels by name"
// @Security BearerAuth
// @Success 200 {object} response.ResponseWithPagination{data=[]hoteldto.ListHotel} "Successfully retrieved list of hotels"
// @Router /hotels [get]
func (hh *HotelHandler) ListHotels(c *gin.Context) {
	ctx := c.Request.Context()

	var req hoteldto.ListHotelRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := hh.hotelUsecase.ListHotels(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching hotels:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get list of hotels")
		return
	}

	pagination := &response.Pagination{}
	message := "Successfully retrieved list of hotels"

	var hotels []hoteldto.ListHotel
	if resp != nil {
		hotels = resp.Hotels
		if len(resp.Hotels) == 0 {
			message = "No users found"
		}
		pagination = response.NewPagination(req.Limit, req.Page, int(resp.Total))
	}

	response.SuccessWithPagination(c, hotels, message, pagination)
}
