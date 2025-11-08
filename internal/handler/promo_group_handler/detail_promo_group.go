package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// DetailPromoGroup godoc
// @Summary Get Promo Group Details
// @Description Retrieve detailed information about a specific promo group by its ID.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param id path int true "Promo Group ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=entity.PromoGroup} "Successfully retrieved promo group details"
// @Router /promo-groups/{id} [get]
func (pgh *PromoGroupHandler) DetailPromoGroup(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error(ctx, "Invalid promo group Id in path:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid promo group Id")
		return
	}

	resp, err := pgh.promoGroupUsecase.DetailPromoGroup(ctx, uint(id))
	if err != nil {
		logger.Error(ctx, "Error getting promo group details:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to get promo group details")
		return
	}

	if resp == nil {
		response.Success(c, nil, "Promo group not found")
		return
	}

	response.Success(c, resp, "Successfully retrieved promo group details")
}
