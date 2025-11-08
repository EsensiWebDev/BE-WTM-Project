package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateStatusBanner godoc
// @Summary Update Banner Status
// @Description Update the status of a banner (active/inactive).
// @Tags Banner
// @Accept json
// @Produce json
// @Param status body bannerdto.UpdateStatusBannerRequest true "Status of the banner"
// @Success 200 {object} response.Response "Successfully updated banner status"
// @Security BearerAuth
// @Router /banners/status [post]
func (bh *BannerHandler) UpdateStatusBanner(c *gin.Context) {
	ctx := c.Request.Context()

	var req bannerdto.UpdateStatusBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error binding request:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "Error validating request:", err.Error())
		if ve := utils.ParseValidationErrors(err); ve != nil {
			response.ValidationError(c, ve)
			return
		}

		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := bh.bannerUsecase.UpdateStatusBanner(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating banner status:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error updating banner status")
		return
	}

	response.Success(c, nil, "Successfully updated banner status")
}
