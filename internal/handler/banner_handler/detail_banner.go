package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// DetailBanner godoc
// @Summary Get Banner Details
// @Description Fetch details of a specific banner by its Id.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path string true "Banner Id"
// @Success 200 {object} response.ResponseWithData{data=bannerdto.BannerData} "Banner details fetched successfully"
// @Security BearerAuth
// @Router /banners/{id} [get]
func (bh *BannerHandler) DetailBanner(c *gin.Context) {
	ctx := c.Request.Context()

	var req bannerdto.DetailBannerRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(ctx, "Error binding query parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	banner, err := bh.bannerUsecase.DetailBanner(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error fetching banner details:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch banner details")
		return
	}

	if banner == nil {
		response.Error(c, http.StatusNotFound, "Banner not found")
		return
	}

	response.Success(c, banner.Banner, "Banner details fetched successfully")
}
