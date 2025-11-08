package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// PromoByID godoc
// @Summary Get Promo by Id
// @Description Retrieve a promo by its Id.
// @Tags Promo
// @Accept json
// @Produce json
// @Param id path string true "Promo Id"
// @Success 200 {object} response.ResponseWithData{data=entity.Promo} "Successfully retrieved promo"
// @Security BearerAuth
// @Router /promos/{id} [get]
func (ph *PromoHandler) PromoByID(c *gin.Context) {
	ctx := c.Request.Context()

	promoID := c.Param("id")
	promoIDUint, err := utils.StringToUint(promoID)
	if err != nil {
		logger.Error(ctx, "Error converting promo Id to uint:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid promo Id")
		return
	}

	promo, err := ph.promoUsecase.PromoByID(ctx, promoIDUint)
	if err != nil {
		logger.Error(ctx, "Error fetching promo by Id:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error fetching promo")
		return
	}

	response.Success(c, promo, "Successfully retrieved promo")
}
