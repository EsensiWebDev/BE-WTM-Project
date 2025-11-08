package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AssignPromoToGroup godoc
// @Summary Assign Promo to Promo Group
// @Description Assign a promo to a specified promo group.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param request body promogroupdto.AssignPromoToGroupRequest true "Assign Promo to Promo Group"
// @Success 200 {object} response.Response "Successfully assigned promo to group"
// @Security BearerAuth
// @Router /promo-groups/promo [post]
func (pgh *PromoGroupHandler) AssignPromoToGroup(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.AssignPromoToGroupRequest
	if err := c.ShouldBind(&req); err != nil {
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

	if err := pgh.promoGroupUsecase.AssignPromoToGroup(ctx, &req); err != nil {
		logger.Error(ctx, "Error assigning promo to group:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to assign promo to group")
		return
	}

	response.Success(c, nil, "Successfully assigned promo to group")
}
