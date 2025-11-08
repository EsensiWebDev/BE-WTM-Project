package promo_group_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
)

// RemovePromoGroup godoc
// @Summary Remove Promo Group
// @Description Remove a promo group by its ID.
// @Tags Promo Group
// @Accept json
// @Produce json
// @Param id path int true "Id of the promo group to be removed"
// @Success 200 {object} response.Response "Successfully removed promo group"
// @Security BearerAuth
// @Router /promo-groups/{id} [delete]
func (pgh *PromoGroupHandler) RemovePromoGroup(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logger.Error(ctx, "Invalid promo group Id in path:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid promo group Id")
		return
	}

	if err := pgh.promoGroupUsecase.RemovePromoGroup(ctx, uint(id)); err != nil {
		logger.Error(ctx, "Error removing promo group:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Failed to remove promo group")
		return
	}

	response.Success(c, nil, "Successfully removed promo group")

}
