package banner_handler

import (
	"fmt"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// UpdateBanner godoc
// @Summary Update Banner
// @Description Update an existing banner with the specified details.
// @Tags Banner
// @Accept json
// @Produce json
// @Param id path string true "Banner Id"
// @Param title formData string false "Title of the banner"
// @Param description formData string false "Description of the banner"
// @Param image formData file false "Image file for the banner"
// @Success 200 {object} response.Response "Successfully updated banner"
// @Security BearerAuth
// @Router /banners/{id} [put]
func (bh *BannerHandler) UpdateBanner(c *gin.Context) {
	ctx := c.Request.Context()

	var reqID bannerdto.DetailBannerRequest
	if err := c.ShouldBindUri(&reqID); err != nil {
		logger.Error(ctx, "Error binding query parameters:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	var req bannerdto.UpsertBannerRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := bh.bannerUsecase.UpsertBanner(ctx, &req, &reqID)
	if err != nil {
		logger.Error(ctx, "Error updating banner:", err.Error())
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed updating banner: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully updated banner")
}
