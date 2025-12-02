package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
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

	var req bannerdto.DetailBannerRequest
	if err := c.ShouldBindUri(&req); err != nil {
		logger.Error(ctx, "Error binding query parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if err := bh.bannerUsecase.RemoveBanner(ctx, &req); err != nil {
		logger.Error(ctx, "Error removing banner:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error removing banner")
		return
	}

	response.Success(c, nil, "Successfully removed banner")
}
