package promo_handler

import (
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// PromoByID godoc
// @Summary Get Promo by Id
// @Description Retrieve a promo by its Id.
// @Tags Promo
// @Accept json
// @Produce json
// @Param id path string true "Promo Id"
// @Success 200 {object} response.ResponseWithData{data=entity.PromoWithExternalID} "Successfully retrieved promo"
// @Security BearerAuth
// @Router /promos/{id} [get]
func (ph *PromoHandler) PromoByID(c *gin.Context) {
	ctx := c.Request.Context()

	promoID := c.Param("id")

	promo, err := ph.promoUsecase.PromoByID(ctx, promoID)
	if err != nil {
		logger.Error(ctx, "Error fetching promo by Id:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error fetching promo")
		return
	}

	response.Success(c, promo, "Successfully retrieved promo")
}
