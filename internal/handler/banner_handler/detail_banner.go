package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// DetailBanner godoc
// @Summary Get Banner Details
// @Description Fetch details of a specific banner by its Id.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path string true "Banner Id"
// @Success 200 {object} response.ResponseWithData{data=entity.Banner} "Banner details fetched successfully"
// @Security BearerAuth
// @Router /banners/{id} [get]
func (bh *BannerHandler) DetailBanner(c *gin.Context) {
	ctx := c.Request.Context()

	bannerID := c.Param("id")
	if bannerID == "" {
		logger.Error(ctx, "Banner Id is required")
		response.Error(c, http.StatusBadRequest, "Banner Id is required")
		return
	}

	bannerIDUint, err := utils.StringToUint(bannerID)
	if err != nil {
		logger.Error(ctx, "Invalid banner Id format:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid banner Id format")
		return
	}

	banner, err := bh.bannerUsecase.DetailBanner(ctx, bannerIDUint)
	if err != nil {
		logger.Error(ctx, "Error fetching banner details:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to fetch banner details")
		return
	}

	if banner == nil {
		response.Error(c, http.StatusNotFound, "Banner not found")
		return
	}

	response.Success(c, banner, "Banner details fetched successfully")
}
