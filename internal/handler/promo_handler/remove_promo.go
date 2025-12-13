package promo_handler

import (
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RemovePromo godoc
// @Summary Remove Promo
// @Description Remove a member from a promo by Id.
// @Tags Promo
// @Accept json
// @Produce json
// @Param id path string true "Promo Id"
// @Success 200 {object} response.Response "Successfully removed promo"
// @Security BearerAuth
// @Router /promos/{id} [delete]
func (ph *PromoHandler) RemovePromo(c *gin.Context) {
	ctx := c.Request.Context()

	promoID := c.Param("id")
	if promoID == "" {
		logger.Error(ctx, "Promo Id is required")
		response.Error(c, http.StatusBadRequest, "Promo Id is required")
		return
	}

	if err := ph.promoUsecase.RemovePromo(ctx, promoID); err != nil {
		logger.Error(ctx, "Error removing promo:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error removing promo")
		return
	}

	response.Success(c, nil, "Successfully removed promo")
}
