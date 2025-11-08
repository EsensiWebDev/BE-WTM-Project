package banner_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/bannerdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdateOrderBanner godoc
// @Summary Update Order Banner
// @Description Update the order of a banner in the promotional banners.
// @Tags Banner
// @Accept json
// @Produce json
// @Param request body bannerdto.UpdateOrderBannerRequest true "Update Order Banner Request"
// @Success 200 {object} response.Response "Successfully updated order banner"
// @Security BearerAuth
// @Router /banners/order [post]
func (bh *BannerHandler) UpdateOrderBanner(c *gin.Context) {
	ctx := c.Request.Context()

	var req bannerdto.UpdateOrderBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Error parsing request:", err.Error())
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

	if err := bh.bannerUsecase.UpdateOrderBanner(ctx, &req); err != nil {
		logger.Error(ctx, "Error updating order banner:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error updating order banner")
		return
	}

	response.Success(c, nil, "Successfully updated order banner")
}
