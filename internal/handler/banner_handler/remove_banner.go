package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemoveBanner godoc
// @Summary Remove Banner
// @Description Remove a banner by Id.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path string true "Banner Id"
// @Success 200 {object} response.Response "Successfully removed banner"
// @Security BearerAuth
// @Router /banners/{id} [delete]
func (bh *BannerHandler) RemoveBanner(c *gin.Context) {
	ctx := c.Request.Context()

	bannerId := c.Param("id")
	if bannerId == "" {
		logger.Error(ctx, "banner id required")
		response.Error(c, http.StatusBadRequest, "Banner Id is required")
		return
	}

	bannerIDUint, err := utils.StringToUint(bannerId)
	if err != nil {
		logger.Error(ctx, "Invalid Banner Id format:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid Banner Id format")
		return
	}

	if err := bh.bannerUsecase.RemoveBanner(ctx, bannerIDUint); err != nil {
		logger.Error(ctx, "Error removing banner:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error removing banner")
		return
	}

	response.Success(c, nil, "Successfully removed banner")
}
