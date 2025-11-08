package promo_handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/response"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

// UpdatePromo godoc
// @Summary Update Promo
// @Description Update an existing promo with the specified details.
// @Tags Promo
// @Accept json
// @Produce json
// @Param id path uint true "Promo ID"
// @Param request body promodto.UpsertPromoRequest true "Promo details"
// @Success 200 {object} response.Response "Successfully updated promo"
// @Security BearerAuth
// @Router /promos/{id} [put]
func (ph *PromoHandler) UpdatePromo(c *gin.Context) {
	ctx := c.Request.Context()

	promoID := c.Param("id")
	if promoID == "" {
		logger.Error(ctx, "Promo Id is required")
		response.Error(c, http.StatusBadRequest, "Promo Id is required")
		return
	}

	promoIDUint, err := utils.StringToUint(promoID)
	if err != nil {
		logger.Error(ctx, "Invalid Promo Id format:", err.Error())
		response.Error(c, http.StatusBadRequest, "Invalid Promo Id format")
		return
	}

	var req *promodto.UpsertPromoRequest
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

	err = ph.promoUsecase.UpsertPromo(ctx, req, &promoIDUint)
	if err != nil {
		logger.Error(ctx, "Error updating promo:", err.Error())
		response.Error(c, http.StatusInternalServerError, "Error updating promo")
		return
	}

	response.Success(c, nil, "Successfully updated promo")
}
