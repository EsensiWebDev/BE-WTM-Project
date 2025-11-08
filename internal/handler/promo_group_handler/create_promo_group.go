package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// CreatePromoGroup godoc
// @Summary Create a new promo group
// @Description Create a new promo group by providing the group name in JSON format.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param request body promogroupdto.CreatePromoGroupRequest true "Promo group creation payload"
// @Security BearerAuth
// @Success 200 {object} response.Response "Successfully created promo group"
// @Router /promo-groups [post]
func (pgh *PromoGroupHandler) CreatePromoGroup(c *gin.Context) {
	ctx := c.Request.Context()
	var req promogroupdto.CreatePromoGroupRequest

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

	err := pgh.promoGroupUsecase.CreatePromoGroup(ctx, req.Name)
	if err != nil {
		logger.Error(ctx, "Error creating promodto group", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to create promodto group")
		return
	}

	response.Success(c, nil, "Successfully created promo group")
}
