package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// RemovePromoGroupMember godoc
// @Summary Remove a member from a promo group
// @Description Remove a user from a specific promo group by providing the required data in JSON format.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param request body promogroupdto.RemovePromoGroupMemberRequest true "Payload to remove member from promo group"
// @Success 200 {object} response.Response "Successfully removed member from promo group"
// @Security BearerAuth
// @Router /promo-groups/members [delete]
func (pgh *PromoGroupHandler) RemovePromoGroupMember(c *gin.Context) {
	ctx := c.Request.Context()
	var req promogroupdto.RemovePromoGroupMemberRequest

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

	err := pgh.promoGroupUsecase.RemovePromoGroupMember(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error removing promo group member", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove promo group member")
		return
	}

	response.Success(c, nil, "Successfully removed member from promo group")
}
