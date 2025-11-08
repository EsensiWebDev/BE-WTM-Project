package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// ListPromoGroupMembers godoc
// @Summary Get members of a promo group
// @Description Retrieve a paginated list of members belonging to a specific promo group using query parameters.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param promo_group_id query int true "Id of the promo group"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} response.ResponseWithPagination{data=[]promogroupdto.ListPromoGroupMemberData} "Successfully retrieved promo group members"
// @Security BearerAuth
// @Router /promo-groups/members [get]
func (pgh *PromoGroupHandler) ListPromoGroupMembers(c *gin.Context) {
	ctx := c.Request.Context()

	var req promogroupdto.ListPromoGroupMemberRequest
	if err := c.ShouldBindQuery(&req); err != nil {
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

	resp, total, err := pgh.promoGroupUsecase.ListPromoGroupMembers(ctx, &req)
	if err != nil {
		logger.Error(ctx, "Error getting promo group members:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get promo group members")
		return
	}

	pagination := response.NewPagination(req.Limit, req.Page, int(total))
	message := "Successfully retrieved promo group members"

	var members []promogroupdto.ListPromoGroupMemberData
	if resp != nil {
		members = resp.PromoGroupMembers
		if len(members) == 0 {
			message = "No promo group members found"
		}
	}

	response.SuccessWithPagination(c, members, message, pagination)
}
