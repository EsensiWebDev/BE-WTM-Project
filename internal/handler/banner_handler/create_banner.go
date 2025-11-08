package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// CreateBanner godoc
// @Summary Create Banner
// @Description Create a new banner with the specified details.
// @Tags Banner
// @Accept json
// @Produce json
// @Param title formData string true "Title of the banner"
// @Param description formData string false "Description of the banner"
// @Param image formData file true "Image file for the banner"
// @Success 200 {object} response.Response "Successfully created banner"
// @Security BearerAuth
// @Router /banners [post]
func (bh *BannerHandler) CreateBanner(c *gin.Context) {
	ctx := c.Request.Context()

	var req bannerdto.UpsertBannerRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Error(ctx, "Error binding request", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := req.ValidateCreate(); err != nil {
		logger.Error(ctx, "Error validating request:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}

		// fallback: unknown validation error
		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err := bh.bannerUsecase.UpsertBanner(ctx, &req, nil)
	if err != nil {
		logger.Error(ctx, "Error creating banner:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to create banner")
		return
	}

	response.Success(c, nil, "Successfully created banner")
}
