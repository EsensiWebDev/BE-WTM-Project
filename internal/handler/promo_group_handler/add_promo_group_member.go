package promo_group_handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// AssignPromoGroupMember godoc
// @Summary Assign users to a promo group
// @Description Assign one or more users to a promo group. This can be used to add new members or update existing users' promo group association.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param request body promogroupdto.AssignPromoGroupMemberRequest true "Payload to add member to promo group"
// @Success 200 {object} response.Response "Successfully assigned member to promo group"
// @Security BearerAuth
// @Router /promo-groups/members [post]
func (pgh *PromoGroupHandler) AssignPromoGroupMember(c *gin.Context) {
	ctx := c.Request.Context()
	var req promogroupdto.AssignPromoGroupMemberRequest

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

	err := pgh.promoGroupUsecase.AssignPromoGroupMember(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error assigning promo group member", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to assign promo group member: %s", err.Error()))
		return
	}

	response.Success(c, nil, "Successfully assigned member to promo group")
}
