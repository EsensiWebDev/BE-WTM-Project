package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemovePromoFromGroup godoc
// @Summary Remove Promo from Promo Group
// @Description Remove a promo from a specified promo group.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param request body promogroupdto.RemovePromoFromGroupRequest true "Remove Promo from Promo Group"
// @Success 200 {object} response.Response "Successfully removed promo from group"
// @Security BearerAuth
// @Router /promo-groups/promo [delete]
func (pgh *PromoGroupHandler) RemovePromoFromGroup(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.RemovePromoFromGroupRequest
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

		// fallback: unknown validation error
		logger.Error(ctx, "Unexpected validation error", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if err := pgh.promoGroupUsecase.RemovePromoFromGroup(ctx, &req); err != nil {
		logger.Error(ctx, "Error removing promo from group", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error removing promo from group")
		return
	}

	response.Success(c, nil, "Successfully removed promo from group")
}
