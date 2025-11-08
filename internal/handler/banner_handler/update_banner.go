package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
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

	bannerID := c.Param("id")
	if bannerID == "" {
		logger.Error(ctx, "Banner Id is required")
		response.Error(c, http.StatusBadRequest, "Banner Id is required")
	}

	bannerIDUint, err := utils.StringToUint(bannerID)
	if err != nil {
		logger.Error(ctx, "Invalid banner Id:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid banner Id format")
	}

	var req *bannerdto.UpsertBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = bh.bannerUsecase.UpsertBanner(ctx, req, &bannerIDUint)
	if err != nil {
		logger.Error(ctx, "Error updating banner:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error updating banner")
		return
	}

	response.Success(c, nil, "Successfully updated banner")
}
